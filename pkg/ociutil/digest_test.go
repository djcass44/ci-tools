package ociutil

import (
	"github.com/Snakdy/container-build-engine/pkg/oci/auth"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDigest(t *testing.T) {
	assert.NotEmpty(t, GetDigest("harbor.dcas.dev/docker.io/library/busybox:latest", auth.Auth{}))
}

func TestSnapshotImage(t *testing.T) {
	t.Run("digest is left alone", func(t *testing.T) {
		in := "busybox:latest@sha256:c118f538365369207c12e5794c3cbfb7b042d950af590ae6c287ede74f29b7d4"
		out, err := SnapshotImage(in, auth.Auth{})
		t.Log(out)
		assert.NoError(t, err)
		assert.EqualValues(t, in, out)
	})
	t.Run("tag is resolved to a digest", func(t *testing.T) {
		in := "busybox"
		out, err := SnapshotImage(in, auth.Auth{})
		t.Log(out)
		assert.NoError(t, err)
		assert.Contains(t, out, "sha256:")
	})
	t.Run("private image is not resolved", func(t *testing.T) {
		out, err := SnapshotImage("harbor.dcas.dev/not-found/not-found:v1.2.3", auth.Auth{})
		assert.Error(t, err)
		assert.Empty(t, out)
	})
}
