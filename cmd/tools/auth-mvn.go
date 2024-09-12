package tools

import (
	"github.com/spf13/cobra"
	"gitlab.com/av1o/mvn-settings-gen/pkg/maven"
	"os"
	"path/filepath"
	"strings"
)

var mvnAuthCmd = &cobra.Command{
	Use:   "mvn-auth",
	Short: "generate a settings.xml for use by Maven",
	RunE:  mvnAuth,
}

const (
	flagMirror    = "mirror"
	flagServer    = "server"
	flagRepo      = "repo"
	flagLocalRepo = "local-repo"

	defaultLocalRepo = "${user.home}/.m2/repository"
)

func init() {
	mvnAuthCmd.Flags().StringArray(flagMirror, nil, "mirror string in the format 'id=name=url=of'")
	mvnAuthCmd.Flags().StringArray(flagServer, nil, "server string in the format 'type=id=username=password'")
	mvnAuthCmd.Flags().StringArray(flagRepo, nil, "repository string in the format 'id=name=url=releases=snapshots'")
	mvnAuthCmd.Flags().String(flagLocalRepo, defaultLocalRepo, "path to the local repository")

	mvnAuthCmd.Flags().String(flagOutputFile, "settings.xml", "path to write the settings.xml file")

	_ = mvnAuthCmd.MarkFlagRequired(flagOutputFile)
}

func mvnAuth(cmd *cobra.Command, _ []string) error {
	mirrors, _ := cmd.Flags().GetStringArray(flagMirror)
	servers, _ := cmd.Flags().GetStringArray(flagServer)
	repos, _ := cmd.Flags().GetStringArray(flagRepo)
	localRepo, _ := cmd.Flags().GetString(flagLocalRepo)
	output, _ := cmd.Flags().GetString(flagOutputFile)

	return generate(output, localRepo, servers, repos, mirrors)
}

func generate(output, localRepo string, servers, repos, mirrors []string) error {
	if strings.TrimSpace(localRepo) == "" {
		localRepo = defaultLocalRepo
	}

	// generate the settings.xml structure
	settings := maven.NewSettings()
	settings.SimpleConfigurer(localRepo, false, false)
	settings.ServerConfigurer(servers)
	settings.MirrorConfigurer(mirrors)
	settings.RepoConfigurer(repos)

	str, err := settings.ToString()
	if err != nil {
		return err
	}

	// write the data to file
	if err := os.WriteFile(filepath.Clean(output), []byte(str), 0644); err != nil {
		return err
	}

	return nil
}
