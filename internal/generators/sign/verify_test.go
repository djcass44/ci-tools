package sign

import (
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerify(t *testing.T) {
	t.Run("valid sig returns no error", func(t *testing.T) {
		err := Verify(&civ1.BuildContext{}, "gcr.io/distroless/static-debian11:nonroot", "./testdata/distroless.pub", true)
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
	err := VerifyAny(&civ1.BuildContext{}, "gcr.io/kaniko-project/executor:debug", "./testdata", true)
	assert.NoError(t, err)
}
