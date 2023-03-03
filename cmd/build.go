package cmd

import "github.com/spf13/cobra"

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "build an application",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

const (
	flagArchetype = "archetype"
)

func init() {
	buildCmd.Flags().StringP(flagArchetype, "a", "", "application recipe to use")

	// flag options
	_ = buildCmd.MarkFlagRequired(flagArchetype)
}
