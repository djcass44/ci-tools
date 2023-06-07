package slsa

import civ1 "github.com/djcass44/ci-tools/internal/api/v1"

const (
	digestSha1   = "sha1"
	digestSha256 = "sha256"

	DefaultBuildType    = "https://github.com/djcass44/ci-tools@v1"
	DefaultVerifierType = "https://github.com/djcass44/ci-tools@v1"

	outProvenance = "provenance.slsa.json"
	outVSA        = "vsa.slsa.json"
	outBuild      = "build.txt"
)

type ExecuteFunc = func(ctx *civ1.BuildContext, r *civ1.BuildRecipe, digest string) error
