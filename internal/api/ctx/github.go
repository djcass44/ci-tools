package ctx

import (
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/pkg/purl"
	"path/filepath"
	"strings"
	"time"
)

func (c *GitHubContext) Normalise() v1.BuildContext {
	imagePath := "ghcr.io/" + c.ProjectNamespace
	// support mono-repos via the PROJECT_PATH env var
	if c.ProjectPath != "" {
		if c.ProjectPathOverride != "" {
			imagePath = filepath.Join(imagePath, strings.ReplaceAll(c.ProjectPathOverride, "/", "-"))
		} else {
			imagePath = filepath.Join(imagePath, strings.ReplaceAll(c.ProjectPath, "/", "-"))
		}
	}
	// collect tags
	tags := []string{
		c.CommitSha,
		"latest",
	}
	if c.CommitRefName != "" {
		tags = append(tags, strings.ReplaceAll(c.CommitRefName, "/", "-"))
	}
	projectUrl := c.ServerURL + "/" + c.ProjectNamespace
	return v1.BuildContext{
		BuildID:    c.JobID,
		Root:       c.ProjectDir,
		Context:    c.ProjectPath,
		ConfigPath: "",
		Image: v1.ImageConfig{
			Name:     imagePath,
			Base:     "",
			Registry: "ghcr.io",
			Username: c.RegistryUsername,
			Password: c.RegistryPassword,
		},
		Tags:       tags,
		Dockerfile: v1.DockerfileConfig{},
		Repo: v1.BuildRepo{
			URL:       projectUrl,
			CommitSha: c.CommitSha,
			Ref:       c.CommitRefName,
			Trunk:     c.CommitRefType != "tag",
		},
		StartTime: time.Now().String(),
		Provider:  purl.TypeGitHub,
	}
}
