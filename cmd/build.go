package cmd

import (
	"fmt"
	"github.com/djcass44/ci-tools/internal/api/ctx"
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/internal/generators/cache"
	"github.com/djcass44/ci-tools/internal/generators/runtime"
	"github.com/djcass44/ci-tools/internal/generators/sbom"
	"github.com/djcass44/ci-tools/internal/generators/sign"
	"github.com/djcass44/ci-tools/internal/generators/slsa"
	"github.com/djcass44/ci-tools/pkg/ociutil"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build an application",
	RunE:  build,
}

const (
	flagRecipe              = "recipe"
	flagRecipeTemplate      = "recipe-template"
	flagRecipeTemplateExtra = "extra-recipe-template"

	flagSkipDockerCFG    = "skip-docker-cfg"
	flagSkipSBOM         = "skip-sbom"
	flagSkipSLSA         = "skip-slsa"
	flagSkipCosignVerify = "skip-cosign-verify"

	flagCosignPublicKey    = "cosign-verify-key"
	flagCosignPublicKeyDir = "cosign-verify-dir"
	flagCosignOffline      = "cosign-offline"
	flagSLSAVersion        = "slsa-version"
)

func init() {
	buildCmd.Flags().StringP(flagRecipe, "a", "", "application recipe to use")
	buildCmd.Flags().String(flagRecipeTemplate, "", "override the default recipe template file")
	buildCmd.Flags().String(flagRecipeTemplateExtra, "", "additional recipe templates to merge with the default recipe template file")

	buildCmd.Flags().Bool(flagSkipDockerCFG, false, "skip generating the registry credentials file even if requested by a recipe")
	buildCmd.Flags().Bool(flagSkipSBOM, false, "skip generating the SBOM")
	buildCmd.Flags().Bool(flagSkipSLSA, false, "skip generating SLSA provenance")
	buildCmd.Flags().Bool(flagSkipCosignVerify, false, "skip verifying the parent image")

	buildCmd.Flags().String(flagCosignPublicKey, "", "path to the Cosign public key used for verifying parent images")
	buildCmd.Flags().String(flagCosignPublicKeyDir, "", "path to the directory containing Cosign public keys used for verifying parent images")
	buildCmd.Flags().Bool(flagCosignOffline, true, "stops Cosign from communicating with any online resources (e.g., fulcio, rekor) when verifying images")
	buildCmd.Flags().String(flagSLSAVersion, "0.2", "slsa version (1.0 or 0.2)")

	// flag options
	_ = buildCmd.MarkFlagRequired(flagRecipe)

	buildCmd.MarkFlagsMutuallyExclusive(flagCosignPublicKey, flagCosignPublicKeyDir, flagSkipCosignVerify)
	buildCmd.MarkFlagsMutuallyExclusive(flagRecipeTemplate, flagRecipeTemplateExtra)
}

func build(cmd *cobra.Command, _ []string) error {
	// read flags
	skipDockerCfg, _ := cmd.Flags().GetBool(flagSkipDockerCFG)
	skipSBOM, _ := cmd.Flags().GetBool(flagSkipSBOM)
	skipSLSA, _ := cmd.Flags().GetBool(flagSkipSLSA)
	skipCosignVerify, _ := cmd.Flags().GetBool(flagSkipCosignVerify)
	arch, _ := cmd.Flags().GetString(flagRecipe)
	tpl, _ := cmd.Flags().GetString(flagRecipeTemplate)
	if tpl != "" {
		log.Printf("using custom recipe template: %s", tpl)
	}
	extras, _ := cmd.Flags().GetString(flagRecipeTemplateExtra)
	extraTemplates := append([]string{tpl}, strings.Split(extras, ",")...)

	cosignPub, _ := cmd.Flags().GetString(flagCosignPublicKey)
	cosignPubDir, _ := cmd.Flags().GetString(flagCosignPublicKeyDir)
	cosignOffline, err := cmd.Flags().GetBool(flagCosignOffline)
	if err != nil {
		log.Println("unable to retrieve the value of the --cosign-offline flag")
		return err
	}
	slsaVersion, _ := cmd.Flags().GetString(flagSLSAVersion)

	// figure out what we need to do
	log.Printf("running recipe: %s", arch)

	context, err := ctx.GetContext()
	if err != nil {
		return err
	}
	context.Builder = arch
	// rewrite the parent image reference to
	// use a digest
	if context.Image.Parent != "" {
		digest, err := ociutil.SnapshotImage(context.Image.Parent)
		if err != nil {
			return err
		}
		context.Image.Parent = digest
	}

	cfg, err := v1.ReadConfigurations(context, extraTemplates...)
	if err != nil {
		return err
	}
	recipe, ok := cfg.Build[arch]
	if !ok {
		return fmt.Errorf("unknown recipe: %s", arch)
	}

	// write OCI credentials file
	// but make sure we don't accidentally overwrite it unless
	// we intend to
	if recipe.DockerCFG && !skipDockerCfg && os.Getenv("CI") != "" {
		if err := v1.WriteDockerCFG(context); err != nil {
			log.Printf("failed to write dockercfg: %s", err)
			return err
		}
	}

	// prepare cache directories
	cache.Execute(context)

	// verify the parent image if one has been specified
	if context.Image.Parent != "" && !skipCosignVerify {
		// if an explicit key has been given, use that
		if cosignPub != "" {
			if err := sign.Verify(context, context.Image.Parent, cosignPub, cosignOffline); err != nil {
				log.Print("failed to verify Cosign signature on parent image")
				return err
			}
		} else {
			if err := sign.VerifyAny(context, context.Image.Parent, cosignPubDir, cosignOffline); err != nil {
				log.Print("failed to verify Cosign signature on parent image")
				return err
			}
		}
	}

	// run the command
	if err := runtime.Execute(context, &recipe); err != nil {
		return err
	}

	// generate the SBOM
	digest := ociutil.GetDigest(fmt.Sprintf("%s:%s", context.Image.Name, context.Repo.CommitSha))
	if !skipSBOM {
		if err := sbom.Execute(context, digest); err != nil {
			return err
		}
	}

	if !skipSLSA {
		f := slsa.ExecuteV02
		if slsaVersion == "1.0" {
			f = slsa.ExecuteV1
		}
		if err := f(context, &recipe, digest); err != nil {
			return err
		}
	}

	return nil
}
