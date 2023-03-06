package ctx

import (
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"path/filepath"
	"strings"
)

type GitLabContext struct {
	ProjectURL  string `env:"CI_PROJECT_URL"`
	ProjectDir  string `env:"CI_PROJECT_DIR"`
	ProjectPath string `env:"PROJECT_PATH"`

	JobImage     string `env:"CI_JOB_IMAGE"`
	JobID        string `env:"CI_JOB_ID"`
	JobStartedAt string `env:"CI_JOB_STARTED_AT"`

	Registry         string `env:"CI_REGISTRY"`
	RegistryImage    string `env:"CI_REGISTRY_IMAGE"`
	RegistryUser     string `env:"CI_REGISTRY_USER"`
	RegistryPassword string `env:"CI_REGISTRY_PASSWORD"`

	CommitBranch   string `env:"CI_COMMIT_BRANCH"`
	CommitTag      string `env:"CI_COMMIT_TAG"`
	CommitSha      string `env:"CI_COMMIT_SHA"`
	CommitShortSha string `env:"CI_COMMIT_SHORT_SHA"`
	CommitRefName  string `env:"CI_COMMIT_REF_NAME"`
}

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
	}
}
