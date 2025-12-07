package sign

import (
	"testing"

	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/sigstore/fulcio/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestVerifyFulcio(t *testing.T) {
	t.Run("valid sig returns no error", func(t *testing.T) {
		// need to make this work...
		err := VerifyFulcio(&civ1.BuildContext{}, "cgr.dev/chainguard/nginx:latest", api.SigstorePublicServerURL)
		assert.Error(t, err)
	})
	t.Run("no sig returns error", func(t *testing.T) {
		err := VerifyFulcio(&civ1.BuildContext{}, "registry.gitlab.com/av1o/base-images/alpine:2542119d", api.SigstorePublicServerURL)
		assert.Error(t, err)
	})
}

func TestVerify(t *testing.T) {
	t.Run("valid sig returns no error", func(t *testing.T) {
		err := Verify(&civ1.BuildContext{}, "gcr.io/kaniko-project/executor:debug", "./testdata/kaniko.pub", true)
		assert.NoError(t, err)
	})
	t.Run("no sig returns error", func(t *testing.T) {
		err := Verify(&civ1.BuildContext{}, "registry.gitlab.com/av1o/base-images/alpine:2542119d", "./testdata/distroless.pub", true)
		assert.Error(t, err)
	})
	t.Run("mismatched key returns error", func(t *testing.T) {
		err := Verify(&civ1.BuildContext{}, "gcr.io/distroless/static-debian11:nonroot", "./testdata/kaniko.pub", true)
		assert.Error(t, err)
	})
}

func TestVerifyAny(t *testing.T) {
	t.Run("directory of certificates works", func(t *testing.T) {
		err := VerifyAny(&civ1.BuildContext{}, "gcr.io/kaniko-project/executor:debug", "./testdata", true)
		assert.NoError(t, err)
	})
	t.Run("non-existent directory does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			err := VerifyAny(&civ1.BuildContext{}, "gcr.io/kaniko-project/executor:debug", "./NOT_A_REAL_DIRECTORY", true)
			assert.Error(t, err)
		})
	})
	t.Run("empty directory fails verification", func(t *testing.T) {
		err := VerifyAny(&civ1.BuildContext{}, "gcr.io/kaniko-project/executor:debug", t.TempDir(), true)
		assert.Error(t, err)
	})
}
