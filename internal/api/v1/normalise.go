package v1

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (c *BuildContext) Normalise() {
	// handle some cleanup in case the env-specific checks
	// couldn't find anything
	wd, _ := os.Getwd()
	if c.Root == "" {
		c.Root = wd
	}

	// handle common stuff
	var buildArgs []string
	for _, e := range os.Environ() {
		k, v, ok := strings.Cut(e, "=")
		if !ok {
			continue
		}
		if strings.HasPrefix(k, EnvBuildArgPrefix) {
			buildArgs = append(buildArgs, fmt.Sprintf("%s=%s", strings.TrimPrefix(k, EnvBuildArgPrefix), v))
		}
	}
	c.Dockerfile.Args = buildArgs
	c.Dockerfile.File = DefaultDockerfile
	if val := os.Getenv(EnvBuildDockerfile); val != "" {
		c.Dockerfile.File = val
	}
	c.Go.ImportPath = os.Getenv(EnvBuildGoImportPath)

	// collect fully-qualified tags
	// e.g. foo.bar/foo/bar:v1.2.3
	fqTags := make([]string, len(c.Tags))
	for i := range fqTags {
		fqTags[i] = fmt.Sprintf("%s:%s", c.Image.Name, c.Tags[i])
	}
	c.FQTags = fqTags
	c.Image.Parent = os.Getenv(EnvBuildImageParent)
	// collect cache configuration
	// as it may differ between
	// CI/CD implementations
	cacheEnabled, err := strconv.ParseBool(os.Getenv(EnvBuildCacheEnabled))
	if err != nil {
		cacheEnabled = true
	}
	cachePath := filepath.Join(c.Root, DefaultCacheName)
	if val := os.Getenv(EnvBuildCachePath); val != "" {
		cachePath = val
	}
	c.Cache.Enabled = cacheEnabled
	c.Cache.Path = cachePath
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
