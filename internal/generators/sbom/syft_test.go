package sbom_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/internal/generators/sbom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	t.Run("image with no digest", func(t *testing.T) {
		tmp := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(tmp, ".cache"), 0750))
		err := sbom.Execute(context.TODO(), &civ1.BuildContext{
			Root: tmp,
			Image: civ1.ImageConfig{
				Name:     "harbor.dcas.dev/docker.io/library/ubuntu",
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
	})
	t.Run("image with digest", func(t *testing.T) {
		tmp := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(tmp, ".cache"), 0750))
		err := sbom.Execute(context.TODO(), &civ1.BuildContext{
			Root: tmp,
			Image: civ1.ImageConfig{
				Name:     "harbor.dcas.dev/docker.io/library/ubuntu",
				Username: "",
				Password: "",
			},
			Tags:   []string{"latest"},
			FQTags: []string{"harbor.dcas.dev/docker.io/library/ubuntu:18.04"},
			Cache: civ1.BuildCache{
				Enabled: true,
				Path:    filepath.Join(tmp, ".cache"),
			},
		}, "harbor.dcas.dev/docker.io/library/ubuntu:18.04@sha256:98706f0f213dbd440021993a82d2f70451a73698315370ae8615cc468ac06624", "98706f0f213dbd440021993a82d2f70451a73698315370ae8615cc468ac06624")
		assert.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(tmp, "sbom.cdx.json"))
		assert.NoError(t, err)
		t.Log(string(data))
	})
	t.Run("helm chart", func(t *testing.T) {
		tmp := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(tmp, ".cache"), 0750))
		err := sbom.Execute(context.TODO(), &civ1.BuildContext{
			Root: tmp,
			Image: civ1.ImageConfig{
				Name:     "harbor.dcas.dev/ghcr.io/grafana/helm-charts/grafana:8.5.12",
				Username: "",
				Password: "",
			},
			Tags:   []string{"latest"},
			FQTags: []string{"harbor.dcas.dev/ghcr.io/grafana/helm-charts/grafana:8.5.12"},
			Cache: civ1.BuildCache{
				Enabled: true,
				Path:    filepath.Join(tmp, ".cache"),
			},
		}, "harbor.dcas.dev/ghcr.io/grafana/helm-charts/grafana:8.5.12", "93c94db36aa18ec8e647c0306e950072ec715aa1f7029337a1d82f6eb08ba212")
		assert.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(tmp, "sbom.cdx.json"))
		assert.NoError(t, err)
		t.Log(string(data))
	})
}
