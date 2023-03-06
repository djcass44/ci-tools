package cmd

import (
	"fmt"
	"github.com/djcass44/ci-tools/internal/api/ctx"
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/internal/generators/runtime"
	"github.com/djcass44/ci-tools/internal/generators/sbom"
	"github.com/djcass44/ci-tools/internal/generators/sign"
	"github.com/djcass44/ci-tools/internal/generators/slsa"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build an application",
	RunE:  build,
}

const (
	flagRecipe           = "recipe"
	flagRecipeTemplate   = "recipe-template"
	flagSkipDockerCFG    = "skip-docker-cfg"
	flagSkipSBOM         = "skip-sbom"
	flagSkipSLSA         = "skip-slsa"
	flagSkipCosignVerify = "skip-cosign-verify"

	flagCosignPublicKey = "cosign-verify-key"
)

func init() {
	buildCmd.Flags().StringP(flagRecipe, "a", "", "application recipe to use")
	buildCmd.Flags().String(flagRecipeTemplate, "", "override the default recipe template file")
	buildCmd.Flags().Bool(flagSkipDockerCFG, false, "skip generating the registry credentials file even if requested by a recipe")
	buildCmd.Flags().Bool(flagSkipSBOM, false, "skip generating the SBOM")
	buildCmd.Flags().Bool(flagSkipSLSA, false, "skip generating SLSA provenance")
	buildCmd.Flags().Bool(flagSkipCosignVerify, false, "skip verifying the parent image")

	buildCmd.Flags().String(flagCosignPublicKey, "", "path to the Cosign public key used for verifying parent images")

	// flag options
	_ = buildCmd.MarkFlagRequired(flagRecipe)
	buildCmd.MarkFlagsMutuallyExclusive(flagCosignPublicKey, flagSkipCosignVerify)
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
	cosignPub, _ := cmd.Flags().GetString(flagCosignPublicKey)

	// figure out what we need to do
	log.Printf("running recipe: %s", arch)

	context, err := ctx.GetContext()
	if err != nil {
		return err
	}
	context.Builder = arch

	cfg, err := v1.ReadConfiguration(tpl, context)
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

	// verify the parent image if one has been specified
	if context.Image.Parent != "" && !skipCosignVerify && cosignPub != "" {
		if err := sign.Verify(context, context.Image.Parent, cosignPub); err != nil {
			log.Print("failed to verify Cosign signature on parent image")
			return err
		}
	}

	// run the command
	if err := runtime.Execute(context, &recipe); err != nil {
		return err
	}

	// generate the SBOM
	if !skipSBOM {
		if err := sbom.Execute(context); err != nil {
			return err
		}
	}

	if !skipSLSA {
		if err := slsa.Execute(context); err != nil {
			return err
		}
	}

	return nil
}
