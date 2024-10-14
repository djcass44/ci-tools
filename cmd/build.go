package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/djcass44/ci-tools/internal/api/ctx"
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/internal/generators/cache"
	"github.com/djcass44/ci-tools/internal/generators/runtime"
	"github.com/djcass44/ci-tools/internal/generators/sbom"
	"github.com/djcass44/ci-tools/internal/generators/sign"
	"github.com/djcass44/ci-tools/internal/generators/slsa"
	"github.com/djcass44/ci-tools/internal/metrics"
	"github.com/djcass44/ci-tools/pkg/in_toto/vsa"
	"github.com/djcass44/ci-tools/pkg/ociutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
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
	flagCosignFulcioURL    = "cosign-fulcio-url"

	flagSLSAVersion       = "slsa-version"
	flagSLSAPredicateOnly = "slsa-predicate-only"

	flagHairpinTag  = "hairpin-tag"
	flagHairpinRepo = "hairpin-repo"
)

const (
	envPrometheusGatewayURL = "PROMETHEUS_PUSH_URL"
	envPrometheusJobName    = "PROMETHEUS_JOB_NAME"
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
	buildCmd.Flags().String(flagCosignFulcioURL, "", "url of the Fulcio instance")
	buildCmd.Flags().Bool(flagCosignOffline, true, "stops Cosign from communicating with any online resources (e.g., fulcio, rekor) when verifying images")

	buildCmd.Flags().String(flagSLSAVersion, vsa.SlsaVersion02, "slsa version (1.0 or 0.2)")
	buildCmd.Flags().Bool(flagSLSAPredicateOnly, false, "do not generate the provenance statement, only the predicate. Needed for compatability with some tools (e.g. cosign)")

	buildCmd.Flags().String(flagHairpinTag, "", "tag to use when validating the image after we've built it. May be needed with some non-standard build tools such as 'helm'. Defaults to the commit sha.")
	buildCmd.Flags().String(flagHairpinRepo, "", "repository to use when validating the image after we've built it. May be needed with some non-standard build tools such as 'helm'.")

	// flag options
	_ = buildCmd.MarkFlagRequired(flagRecipe)

	buildCmd.MarkFlagsMutuallyExclusive(flagCosignPublicKey, flagCosignPublicKeyDir, flagSkipCosignVerify)
	buildCmd.MarkFlagsMutuallyExclusive(flagRecipeTemplate, flagRecipeTemplateExtra)
}

func build(cmd *cobra.Command, _ []string) error {
	totalTimer := prometheus.NewTimer(metrics.MetricOpsDuration)
	defer pushMetrics(cmd.Context())

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
	cosignFulcioURL, _ := cmd.Flags().GetString(flagCosignFulcioURL)
	cosignOffline, err := cmd.Flags().GetBool(flagCosignOffline)
	if err != nil {
		log.Println("unable to retrieve the value of the --cosign-offline flag")
		return err
	}
	slsaVersion, _ := cmd.Flags().GetString(flagSLSAVersion)
	slsaPredicateOnly, _ := cmd.Flags().GetBool(flagSLSAPredicateOnly)

	// figure out what we need to do
	log.Printf("running recipe: %s", arch)

	context, err := ctx.GetContext()
	if err != nil {
		return err
	}
	context.Builder = arch
	auth := context.Auth()
	// rewrite the parent image reference to
	// use a digest
	if context.Image.Parent != "" {
		digest, err := ociutil.SnapshotImage(context.Image.Parent, auth)
		if err != nil {
			log.Print("unable to verify parent image")
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

	defaultLabelValues := []string{arch, context.Repo.URL, context.Provider, context.Context}
	metrics.MetricOpsBuilds.WithLabelValues(defaultLabelValues...).Inc()

	// write OCI credentials file
	// but make sure we don't accidentally overwrite it unless
	// we intend to
	if recipe.DockerCFG && !skipDockerCfg && os.Getenv("CI") != "" {
		metrics.MetricDockerCfgGenerated.WithLabelValues(defaultLabelValues...).Inc()
		if err := v1.WriteDockerCFG(context); err != nil {
			log.Printf("failed to write dockercfg: %s", err)
			return err
		}
	}

	// prepare cache directories
	cache.Execute(context)

	hairpinTag, _ := cmd.Flags().GetString(flagHairpinTag)
	if hairpinTag == "" {
		hairpinTag = context.Repo.CommitSha
	}
	hairpinRepo, _ := cmd.Flags().GetString(flagHairpinRepo)
	if hairpinRepo == "" {
		hairpinRepo = strings.TrimPrefix(context.Image.Name, context.Image.Registry)
	}

	// verify the parent image if one has been specified
	if context.Image.Parent != "" && !skipCosignVerify {
		// if an explicit key has been given, use that
		if cosignPub != "" {
			metrics.MetricOciVerify.WithLabelValues("public_key").Inc()
			if err := sign.Verify(context, context.Image.Parent, cosignPub, cosignOffline); err != nil {
				log.Print("failed to verify Cosign signature on parent image")
				return err
			}
		} else {
			if cosignFulcioURL != "" {
				metrics.MetricOciVerify.WithLabelValues("fulcio").Inc()
				if err := sign.VerifyFulcio(context, context.Image.Parent, cosignFulcioURL); err != nil {
					log.Printf("failed to verify parent image signature using Fulcio")
				}
			}
			metrics.MetricOciVerify.WithLabelValues("any").Inc()
			if err := sign.VerifyAny(context, context.Image.Parent, cosignPubDir, cosignOffline); err != nil {
				log.Print("failed to verify Cosign signature on parent image")
				return err
			}
		}
	}

	// run the command
	executionTimer := prometheus.NewTimer(metrics.MetricOpsBuildDuration)
	if err := runtime.Execute(context, &recipe); err != nil {
		metrics.MetricOpsBuildErrors.WithLabelValues(defaultLabelValues...).Inc()
		return err
	}
	executionTimer.ObserveDuration()

	// generate the SBOM
	provenanceTimer := prometheus.NewTimer(metrics.MetricOpsProvenanceDuration)
	imageRef := getImageRef(context.Image.Registry, hairpinRepo, hairpinTag)
	digest := ociutil.GetDigest(imageRef, auth)
	// catch digest failures and error out early
	// otherwise the SBOM/SLSA generation wigs out
	if digest == "" {
		return errors.New("cannot generate metadata without a valid digest")
	}
	if !skipSBOM {
		metrics.MetricBOMGenerated.Inc()
		if err := sbom.Execute(cmd.Context(), context, imageRef, digest); err != nil {
			return err
		}
	}

	if !skipSLSA {
		f := slsa.ExecuteV02
		if slsaVersion == vsa.SlsaVersion1 {
			f = slsa.ExecuteV1
		}
		metrics.MetricProvenanceGenerated.WithLabelValues(slsaVersion).Inc()
		if err := f(context, &recipe, imageRef, digest, slsaPredicateOnly); err != nil {
			return err
		}
	}
	provenanceTimer.ObserveDuration()
	totalTimer.ObserveDuration()
	metrics.MetricOpsBuildSuccess.WithLabelValues(defaultLabelValues...).Inc()

	return nil
}

func getImageRef(registry, repo, tag string) string {
	return fmt.Sprintf("%s/%s:%s", strings.TrimSuffix(registry, "/"), strings.TrimPrefix(repo, "/"), tag)
}

func pushMetrics(ctx context.Context) {
	// push prometheus metrics
	if os.Getenv(envPrometheusGatewayURL) == "" {
		log.Println("skipping metrics upload")
		return
	}
	log.Println("uploading metrics")
	err := push.New(os.Getenv(envPrometheusGatewayURL), os.Getenv(envPrometheusJobName)).PushContext(ctx)
	if err != nil {
		log.Printf("failed to push prometheus metrics: %v", err)
	}
}
