package cmd

import (
	"github.com/djcass44/ci-tools/cmd/tools"
	"github.com/spf13/cobra"
	"os"
)

var command = &cobra.Command{
	Use:          "ci",
	Short:        "normalises application build tools",
	SilenceUsage: true,
}

func init() {
	command.AddCommand(buildCmd, dockerfileCmd, verifyCmd, tools.Command)
}

func Execute() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
