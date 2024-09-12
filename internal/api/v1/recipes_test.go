package v1_test

import (
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/internal/generators/runtime"
	"github.com/djcass44/ci-tools/pkg/purl"
	"github.com/franiglesias/golden"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestRecipes(t *testing.T) {
	buildContext := &v1.BuildContext{
		Builder: "https://example.com@v1",
		BuildID: "12345",
		Root:    "/builds/foo/bar",
		Context: "",
		Image: v1.ImageConfig{
			Parent:   "registry.gitlab.example.com/foo/base-images/run:latest",
			Base:     "registry.gitlab.example.com/foo/base-images/build:latest",
			Name:     "foo/bar",
			Registry: "registry.gitlab.example.com",
			Username: "gitlab-ci-token",
			Password: "hunter2",
		},
		Tags: []string{
			"latest",
			"main",
			"deadbeef",
			"v1.2.3",
		},
		FQTags: []string{
			"registry.gitlab.example.com/foo/bar:latest",
			"registry.gitlab.example.com/foo/bar:main",
			"registry.gitlab.example.com/foo/bar:deadbeef",
			"registry.gitlab.example.com/foo/bar:v1.2.3",
		},
		Dockerfile: v1.DockerfileConfig{},
		Go:         v1.GoConfig{},
		Repo: v1.BuildRepo{
			URL:       "https://gitlab.example.com/foo/bar.git",
			CommitSha: "deadbeef",
			Ref:       "main",
			Trunk:     true,
		},
		Cache: v1.BuildCache{
			Enabled: true,
			Path:    "/builds/foo/bar/.cache",
		},
		StartTime:  "2021-11-05T20:12:38Z",
		Provider:   purl.TypeGitLab,
		ConfigPath: ".gitlab-ci.yml",
		ExtraArgs:  nil,
	}

	// read the recipe config
	recipes, err := v1.ReadConfiguration("./recipes.tpl.yaml", buildContext)
	assert.NoError(t, err)
	assert.NotNil(t, recipes)

	// prepare the execution plan
	for k, v := range recipes.Build {
		t.Run(k, func(t *testing.T) {
			plan := runtime.GetExecutionPlan(buildContext, &v, false)
			assert.NotNil(t, plan)

			sort.Strings(plan.Args)
			sort.Strings(plan.Env)

			golden.Verify(t, plan)
		})
	}
}
