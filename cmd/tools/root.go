package tools

import "github.com/spf13/cobra"

var Command = &cobra.Command{
	Use:   "tools",
	Short: "additional tools",
}

func init() {
	Command.AddCommand(sarif2GitLabCmd, mvnAuthCmd)
}
