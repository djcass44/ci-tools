package v1_test

import (
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestReadConfiguration(t *testing.T) {
	cfg, err := v1.ReadConfiguration("", &v1.BuildContext{
		Root:    "/foo/bar",
		Context: "samples/python",
		Image: v1.ImageConfig{
			Name: "localhost:5000/myrepo",
		},
		Tags:       nil,
		FQTags:     nil,
		Dockerfile: v1.DockerfileConfig{},
		Cache: v1.BuildCache{
			Enabled: true,
		},
	})
	assert.NoError(t, err)
	t.Logf("%+v", cfg)
}

func TestReadExtraArgs(t *testing.T) {
	require.NoError(t, os.Setenv("BUILD_EXTRA_ARGS", "-e=foo,foo=zoo"))

	bctx := &v1.BuildContext{
		Repo: v1.BuildRepo{
			CommitSha: "deadbeef",
		},
	}
	bctx.Normalise()
	recipes, err := v1.ReadConfiguration("./testdata/recipes-extra-args.tpl.yaml", bctx)
	assert.NoError(t, err)
	assert.Len(t, recipes.Build, 1)

	r, ok := recipes.Build["echo"]
	assert.True(t, ok)
	assert.ElementsMatch(t, []string{"echo", "-e=foo", "foo=zoo"}, r.Args)
}

func TestReadConfigurations(t *testing.T) {
	t.Run("later recipes are preferred", func(t *testing.T) {
		recipes, err := v1.ReadConfigurations(&v1.BuildContext{
			Root:    "/foo/bar",
			Context: "samples/python",
			Image: v1.ImageConfig{
				Name: "localhost:5000/myrepo",
			},
			Tags:       nil,
			FQTags:     nil,
			Dockerfile: v1.DockerfileConfig{},
			Cache: v1.BuildCache{
				Enabled: true,
			},
		}, "testdata/recipes-mock.tpl.yaml", "testdata/recipes-mock2.tpl.yaml")
		assert.NoError(t, err)
		assert.Len(t, recipes.Build["echo"].Args, 2)
	})
	t.Run("defaults are still loaded", func(t *testing.T) {
		recipes, err := v1.ReadConfigurations(&v1.BuildContext{
			Root:    "/foo/bar",
			Context: "samples/python",
			Image: v1.ImageConfig{
				Name: "localhost:5000/myrepo",
			},
			Tags:       nil,
			FQTags:     nil,
			Dockerfile: v1.DockerfileConfig{},
			Cache: v1.BuildCache{
				Enabled: true,
			},
		}, " testdata/recipes-mock.tpl.yaml", "")
		assert.NoError(t, err)
		assert.Greater(t, len(recipes.Build), 1)
	})
}
