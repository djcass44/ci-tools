package cmd

import (
	"fmt"
	"github.com/djcass44/ci-tools/internal/api/ctx"
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/internal/runtime"
	"github.com/spf13/cobra"
	"log"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build an application",
	RunE:  build,
}

const (
	flagArchetype      = "archetype"
	flagRecipeTemplate = "recipe-template"
	flagSkipDockerCFG  = "skip-docker-cfg"
)

func init() {
	buildCmd.Flags().StringP(flagArchetype, "a", "", "application recipe to use")
	buildCmd.Flags().String(flagRecipeTemplate, "", "override the default recipe template file")
	buildCmd.Flags().Bool(flagSkipDockerCFG, false, "skip generating the registry credentials file even if requested by a recipe")

	// flag options
	_ = buildCmd.MarkFlagRequired(flagArchetype)
}

func build(cmd *cobra.Command, _ []string) error {
	// read flags
	skipDockerCfg, _ := cmd.Flags().GetBool(flagSkipDockerCFG)
	arch, _ := cmd.Flags().GetString(flagArchetype)
	tpl, _ := cmd.Flags().GetString(flagRecipeTemplate)
	if tpl != "" {
		log.Printf("using custom recipe template: %s", tpl)
	}

	// figure out what we need to do
	log.Printf("running recipe: %s", arch)

	context, err := ctx.GetContext()
	if err != nil {
		return err
	}

	cfg, err := v1.ReadConfiguration(tpl, &context)
	if err != nil {
		return err
	}
	recipe, ok := cfg.Build[arch]
	if !ok {
		return fmt.Errorf("unknown recipe: %s", arch)
	}

	// write OCI credentials file
	if recipe.DockerCFG && !skipDockerCfg {
		if err := v1.WriteDockerCFG(&context); err != nil {
			log.Printf("failed to write dockercfg: %s", err)
			return err
		}
	}

	// run the command
	return runtime.Execute(&recipe)
}
