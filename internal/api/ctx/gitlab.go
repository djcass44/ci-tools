package ctx

import (
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/pkg/purl"
	"path/filepath"
	"strings"
)

func (c *GitLabContext) Normalise() v1.BuildContext {
	imagePath := c.RegistryImage
	// support mono-repos via the PROJECT_PATH env var
	if c.ProjectPath != "" {
		imagePath = filepath.Join(imagePath, strings.ReplaceAll(c.ProjectPath, "/", "-"))
	}
	// collect tags
	tags := []string{
		c.CommitSha,
		c.CommitShortSha,
		"latest",
	}
	if c.CommitTag != "" {
		tags = append(tags, strings.ReplaceAll(c.CommitTag, "/", "-"))
	}
	if c.CommitBranch != "" {
		tags = append(tags, strings.ReplaceAll(c.CommitBranch, "/", "-"))
	}
	return v1.BuildContext{
		BuildID: c.JobID,
		Root:    c.ProjectDir,
		Context: c.ProjectPath,
		Image: v1.ImageConfig{
			Name:     imagePath,
			Base:     c.JobImage,
			Registry: c.Registry,
			Username: c.RegistryUser,
			Password: c.RegistryPassword,
		},
		Tags:       tags,
		Dockerfile: v1.DockerfileConfig{},
		Repo: v1.BuildRepo{
			URL:       c.ProjectURL,
			CommitSha: c.CommitSha,
			Ref:       c.CommitRefName,
		},
		StartTime: c.JobStartedAt,
		Provider:  purl.TypeGitLab,
	}
}
