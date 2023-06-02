package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var command = &cobra.Command{
	Use:          "ci",
	Short:        "normalises application build tools",
	SilenceUsage: true,
}

func init() {
	command.AddCommand(buildCmd, dockerfileCmd, verifyCmd)
}

func Execute() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
