package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var command = &cobra.Command{
	Use:   "ci",
	Short: "normalises application build tools",
}

func init() {
	command.AddCommand(buildCmd)
}

func Execute() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
