package v1

import (
	"fmt"
	"os"
	"strings"
)

func (c *BuildContext) Normalise() {
	// handle some cleanup in case the env-specific checks
	// couldn't find anything
	wd, _ := os.Getwd()
	if c.Root == "" {
		c.Root = wd
	}
	if c.Context == "" {
		c.Context = c.Root
	}

	// handle common stuff
	var buildArgs []string
	for _, e := range os.Environ() {
		k, v, ok := strings.Cut(e, "=")
		if !ok {
			continue
		}
		if strings.HasPrefix(k, "BUILD_ARG_") {
			buildArgs = append(buildArgs, fmt.Sprintf("%s=%s", strings.TrimPrefix(k, "BUILD_ARG_"), v))
		}
	}
	c.Dockerfile.Args = buildArgs
	c.Dockerfile.File = "Dockerfile"
	if val := os.Getenv("BUILD_DOCKERFILE"); val != "" {
		c.Dockerfile.File = val
	}

	// collect fully-qualified tags
	// e.g. foo.bar/foo/bar:v1.2.3
	fqTags := make([]string, len(c.Tags))
	for i := range fqTags {
		fqTags[i] = fmt.Sprintf("%s:%s", c.Image.Name, c.Tags[i])
	}
	c.FQTags = fqTags
	c.Image.Parent = os.Getenv("BUILD_IMAGE_PARENT")
}

func (c *BuildContext) DockerCFG() string {
	return fmt.Sprintf(`{
	"auths": {
		"%s": {
			"username": "%s",
			"password": "%s"
		}
	}
}`, c.Image.Registry, c.Image.Username, c.Image.Password)
}
