package ctx

import (
	"fmt"
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"path/filepath"
	"strings"
)

type GitLabContext struct {
	ProjectDir  string `envconfig:"CI_PROJECT_DIR"`
	ProjectPath string `envconfig:"PROJECT_PATH"`

	Registry         string `envconfig:"CI_REGISTRY"`
	RegistryImage    string `envconfig:"CI_REGISTRY_IMAGE"`
	RegistryUser     string `envconfig:"CI_REGISTRY_USER"`
	RegistryPassword string `envconfig:"CI_REGISTRY_PASSWORD"`

	CommitBranch   string `envconfig:"CI_COMMIT_BRANCH"`
	CommitTag      string `envconfig:"CI_COMMIT_TAG"`
	CommitSha      string `envconfig:"CI_COMMIT_SHA"`
	CommitShortSha string `envconfig:"CI_COMMIT_SHORT_SHA"`
}

func (c *GitLabContext) Normalise() v1.BuildContext {
	imagePath := c.RegistryImage
	// support monorepos via the PROJECT_PATH env var
	if c.ProjectPath != "" {
		imagePath = filepath.Join(imagePath, c.ProjectPath)
	}
	// collect tags
	tags := []string{
		"latest",
		c.CommitSha,
		c.CommitShortSha,
	}
	if c.CommitTag != "" {
		tags = append(tags, c.CommitTag)
	}
	if c.CommitBranch != "" {
		tags = append(tags, strings.ReplaceAll(c.CommitBranch, "/", "-"))
	}
	// collect fully-qualified tags
	// e.g. foo.bar/foo/bar:v1.2.3
	fqTags := make([]string, len(tags))
	for i := range fqTags {
		fqTags[i] = fmt.Sprintf("%s:%s", c.RegistryImage, tags[i])
	}
	return v1.BuildContext{
		Root:    c.ProjectDir,
		Context: c.ProjectPath,
		Image: v1.ImageConfig{
			Name:     c.RegistryImage,
			Registry: c.Registry,
			Username: c.RegistryUser,
			Password: c.RegistryPassword,
		},
		Tags:       tags,
		FQTags:     fqTags,
		Dockerfile: v1.DockerfileConfig{},
	}
}
