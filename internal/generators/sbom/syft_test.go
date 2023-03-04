package sbom_test

import (
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
	err := sbom.Execute(&civ1.BuildContext{
		Root: tmp,
		Image: civ1.ImageConfig{
			Name:     "harbor.dcas.dev/docker.io/library/busybox",
			Username: "",
			Password: "",
		},
		Tags:   []string{"latest"},
		FQTags: []string{"harbor.dcas.dev/docker.io/library/busybox:latest"},
	})
	assert.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(tmp, "sbom.cdx.json"))
	assert.NoError(t, err)
	t.Log(string(data))
}
