package tools

import (
	"encoding/json"
	"github.com/owenrumney/go-sarif/sarif"
	"github.com/spf13/cobra"
	"gitlab.com/av1o/gitlab-cq/pkg/generic"
	"os"
	"path/filepath"
)

var sarif2GitLabCmd = &cobra.Command{
	Use:   "sarif2gitlab",
	Short: "Convert a SARIF file into a GitLab CodeQuality report",
	Args:  cobra.NoArgs,
	RunE:  sarif2GitLab,
}

const (
	flagInputFile  = "input"
	flagOutputFile = "output"
)

func init() {
	sarif2GitLabCmd.Flags().StringP(flagInputFile, "i", "", "path to the SARIF file")
	sarif2GitLabCmd.Flags().StringP(flagOutputFile, "o", "", "path to write the GitLab CodeQuality report")

	_ = sarif2GitLabCmd.MarkFlagRequired(flagInputFile)
	_ = sarif2GitLabCmd.MarkFlagRequired(flagOutputFile)
}

func sarif2GitLab(cmd *cobra.Command, _ []string) error {
	input, _ := cmd.Flags().GetString(flagInputFile)
	output, _ := cmd.Flags().GetString(flagOutputFile)

	// read the sarif report
	report, err := sarif.Open(filepath.Clean(input))
	if err != nil {
		return err
	}
	// convert it to a gitlab report
	gitlabReport := generic.AsGitLab(report)

	// create the output file
	f, err := os.Create(filepath.Clean(output))
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")

	// write json to the file
	return enc.Encode(&gitlabReport)
}
