package sbom_test

import (
	"context"
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/internal/generators/sbom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestExecute(t *testing.T) {
	tmp := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(tmp, ".cache"), 0750))
	err := sbom.Execute(context.TODO(), &civ1.BuildContext{
		Root: tmp,
		Image: civ1.ImageConfig{
			Name:     " harbor.dcas.dev/docker.io/library/ubuntu",
			Username: "",
			Password: "",
		},
		Tags:   []string{"latest"},
		FQTags: []string{"harbor.dcas.dev/docker.io/library/ubuntu:18.04"},
		Cache: civ1.BuildCache{
			Enabled: true,
			Path:    filepath.Join(tmp, ".cache"),
		},
	}, "harbor.dcas.dev/docker.io/library/ubuntu:18.04", "98706f0f213dbd440021993a82d2f70451a73698315370ae8615cc468ac06624")
	assert.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(tmp, "sbom.cdx.json"))
	assert.NoError(t, err)
	t.Log(string(data))
}
