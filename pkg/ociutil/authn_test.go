package ociutil

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicKeychain_Resolve(t *testing.T) {
	kc := NewBasicKeychain(Auth{
		Registry: "registry.gitlab.com",
		Username: "gitlab-ci-token",
		Password: "CI_JOB_TOKEN",
	})

	registryGitlab, _ := name.NewRegistry("registry.gitlab.com")
	registryGithub, _ := name.NewRegistry("ghcr.io")

	t.Run("gitlab", func(t *testing.T) {
		auth, err := kc.Resolve(registryGitlab)
		assert.NoError(t, err)
		assert.IsType(t, &authn.Basic{}, auth)
	})
	t.Run("github", func(t *testing.T) {
		auth, err := kc.Resolve(registryGithub)
		assert.NoError(t, err)
		assert.IsType(t, authn.Anonymous, auth)
	})
}
