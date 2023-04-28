package cmd

import (
	"fmt"
	"github.com/djcass44/ci-tools/internal/api/ctx"
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/internal/generators/dockerfile"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
	"strings"
)

var dockerfileCmd = &cobra.Command{
	Use:   "dockerfile",
	Short: "retrieve a Dockerfile",
	RunE:  retrieve,
}

const (
	flagName    = "name"
	flagOutPath = "out"
)

func init() {
	dockerfileCmd.Flags().StringP(flagName, "n", "", "Dockerfile to retrieve")
	dockerfileCmd.Flags().StringP(flagOutPath, "o", "", "Directory to save the Dockerfile. Defaults to 'project-dir/monorepo-dir/Dockerfile'")

	dockerfileCmd.Flags().String(flagRecipeTemplate, "", "override the default recipe template file")
	dockerfileCmd.Flags().String(flagRecipeTemplateExtra, "", "additional recipe templates to merge with the default recipe template file")

	// flag options
	_ = dockerfileCmd.MarkFlagRequired(flagName)
	dockerfileCmd.MarkFlagsMutuallyExclusive(flagRecipeTemplate, flagRecipeTemplateExtra)
}

func retrieve(cmd *cobra.Command, _ []string) error {
	// read flags
	name, _ := cmd.Flags().GetString(flagName)
	out, _ := cmd.Flags().GetString(flagOutPath)

	tpl, _ := cmd.Flags().GetString(flagRecipeTemplate)
	if tpl != "" {
		log.Printf("using custom recipe template: %s", tpl)
	}
	extras, _ := cmd.Flags().GetString(flagRecipeTemplateExtra)
	extraTemplates := append([]string{tpl}, strings.Split(extras, ",")...)

	log.Printf("retrieving Dockerfile: %s", name)

	context, err := ctx.GetContext()
	if err != nil {
		return err
	}

	cfg, err := v1.ReadConfigurations(context, extraTemplates...)
	if err != nil {
		return err
	}
	df, ok := cfg.Dockerfiles[name]
	if !ok {
		return fmt.Errorf("failed to locate Dockerfile: %s", name)
	}
	// if the user hasn't specified an output location,
	// then we figure it out based on the environmental
	// context
	if out == "" {
		out = filepath.Join(context.Root, context.Context, "Dockerfile")
	}
	// retrieve the Dockerfile and write
	// it to disk
	if err := dockerfile.Get(&df.Content, out); err != nil {
		return err
	}
	log.Printf("successfully wrote Dockerfile to: %s", out)

	return nil
}
