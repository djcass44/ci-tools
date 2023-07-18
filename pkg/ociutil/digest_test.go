package ociutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDigest(t *testing.T) {
	assert.NotEmpty(t, GetDigest("harbor.dcas.dev/docker.io/library/busybox:latest", Auth{}))
}

func TestSnapshotImage(t *testing.T) {
	t.Run("digest is left alone", func(t *testing.T) {
		in := "busybox:latest@sha256:c118f538365369207c12e5794c3cbfb7b042d950af590ae6c287ede74f29b7d4"
		out, err := SnapshotImage(in, Auth{})
		t.Log(out)
		assert.NoError(t, err)
		assert.EqualValues(t, in, out)
	})
	t.Run("tag is resolved to a digest", func(t *testing.T) {
		in := "busybox"
		out, err := SnapshotImage(in, Auth{})
		t.Log(out)
		assert.NoError(t, err)
		assert.Contains(t, out, "sha256:")
	})
}
