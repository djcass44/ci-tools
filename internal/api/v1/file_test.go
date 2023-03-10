package v1_test

import (
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/stretchr/testify/assert"
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
