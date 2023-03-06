package slsa_test

import (
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/internal/generators/slsa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestExecute(t *testing.T) {
	tmp := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(tmp, ".cache"), 0750))
	err := slsa.Execute(&civ1.BuildContext{
		Root: tmp,
		Image: civ1.ImageConfig{
			Name:   "harbor.dcas.dev/docker.io/library/busybox",
			Parent: "harbor.dcas.dev/docker.io/library/busybox",
			Base:   "harbor.dcas.dev/docker.io/library/busybox",
		},
		Tags:    []string{"latest"},
		FQTags:  []string{"harbor.dcas.dev/docker.io/library/busybox:latest"},
		BuildID: "1",
		Builder: "foo.bar",
		Repo: civ1.BuildRepo{
			URL:       "https://example.org",
			CommitSha: "deadbeef",
			Ref:       "v1.2.3",
		},
		StartTime: "2021-11-05T20:12:38Z",
	})
	assert.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(tmp, "provenance.slsa.json"))
	assert.NoError(t, err)
	t.Log(string(data))
}
