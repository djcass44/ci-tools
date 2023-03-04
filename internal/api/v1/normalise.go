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
}
