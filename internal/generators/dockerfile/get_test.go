package dockerfile

import (
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestGetSuccess(t *testing.T) {
	var cases = []struct {
		name string
		in   *v1.DockerfileContent
	}{
		{
			"file",
			&v1.DockerfileContent{File: "./testdata/Dockerfile"},
		},
		{
			"inline",
			&v1.DockerfileContent{Inline: "FROM harbor.dcas.dev/docker.io/library/busybox:latest"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "Dockerfile")
			err := Get(tt.in, out)
			assert.NoError(t, err)
			assert.FileExists(t, out)
		})
	}
}

func TestGetFail(t *testing.T) {
	var cases = []struct {
		name string
		in   *v1.DockerfileContent
	}{
		{
			"file",
			&v1.DockerfileContent{File: "./testdata/notreal.Dockerfile"},
		},
		{
			"inline",
			&v1.DockerfileContent{},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "Dockerfile")
			err := Get(tt.in, out)
			assert.Error(t, err)
			assert.NoFileExists(t, out)
		})
	}
}
