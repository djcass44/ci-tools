package ctx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGitLabContext_Normalise(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctx := &GitLabContext{
			CommitBranch:   "main",
			CommitSha:      "deadbeef",
			CommitShortSha: "deadbee",
			Registry:       "registry.gitlab.com",
			RegistryImage:  "registry.gitlab.com/foo/bar",
		}
		out := ctx.Normalise()
		assert.EqualValues(t, "registry.gitlab.com/foo/bar", out.Image.Name)
		assert.Empty(t, out.Context)
	})
	t.Run("project path", func(t *testing.T) {
		ctx := &GitLabContext{
			ProjectPath:    "app1",
			CommitBranch:   "main",
			CommitSha:      "deadbeef",
			CommitShortSha: "deadbee",
			Registry:       "registry.gitlab.com",
			RegistryImage:  "registry.gitlab.com/foo/bar",
		}
		out := ctx.Normalise()
		assert.EqualValues(t, "registry.gitlab.com/foo/bar/app1", out.Image.Name)
		assert.EqualValues(t, "app1", out.Context)
	})
	t.Run("project path and override", func(t *testing.T) {
		ctx := &GitLabContext{
			ProjectPath:         "app1",
			ProjectPathOverride: "foobar",
			CommitBranch:        "main",
			CommitSha:           "deadbeef",
			CommitShortSha:      "deadbee",
			Registry:            "registry.gitlab.com",
			RegistryImage:       "registry.gitlab.com/foo/bar",
		}
		out := ctx.Normalise()
		assert.EqualValues(t, "registry.gitlab.com/foo/bar/foobar", out.Image.Name)
		assert.EqualValues(t, "app1", out.Context)
	})
}
