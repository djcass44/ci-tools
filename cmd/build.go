package cmd

import (
	"fmt"
	"github.com/djcass44/ci-tools/internal/api/ctx"
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
	"log"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build an application",
	RunE:  build,
}

const (
	flagArchetype = "archetype"
)

func init() {
	buildCmd.Flags().StringP(flagArchetype, "a", "", "application recipe to use")

	// flag options
	_ = buildCmd.MarkFlagRequired(flagArchetype)
}

func build(cmd *cobra.Command, _ []string) error {
	arch, _ := cmd.Flags().GetString(flagArchetype)
	log.Printf("running recipe: %s", arch)

	var context ctx.GitLabContext
	if err := envconfig.Process("", &context); err != nil {
		return err
	}
	bc := context.Normalise()
	bc.Normalise()

	cfg, err := v1.ReadConfiguration("recipes.tpl.yaml", &bc)
	if err != nil {
		return err
	}
	recipe, ok := cfg.Build[arch]
	if !ok {
		return fmt.Errorf("unknown recipe: %s", arch)
	}
	log.Printf("%+v", recipe)

	return nil
}
