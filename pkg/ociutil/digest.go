package ociutil

import (
	"errors"
	"fmt"
	"github.com/Snakdy/container-build-engine/pkg/oci/auth"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"log"
	"strings"
)

// GetDigest returns the SHA256 digest of
// a given container image.
//
// If the digest cannot be found, an
// empty string is returned.
func GetDigest(target string, authn auth.Auth) string {
	sha, err := crane.Digest(target, crane.WithAuthFromKeychain(auth.KeyChain(authn)))
	if err != nil {
		log.Printf("unable to generate digest for: %s (%s)", target, err)
		return ""
	}
	return strings.TrimPrefix(sha, "sha256:")
}

// SnapshotImage mangles a given OCI string so that
// there is a digest present.
func SnapshotImage(target string, authn auth.Auth) (string, error) {
	ref, err := name.ParseReference(target)
	if err != nil {
		return "", err
	}
	// if the name already contains a digest, we can
	// just return straight away
	if d, err := name.NewDigest(ref.String()); err == nil {
		return d.String(), nil
	}
	// if there's no tag and no digest, throw an error
	if _, err := name.NewTag(ref.String()); err != nil {
		return "", err
	}
	// otherwise grab the digest for the tag and
	// splice it in
	digest := GetDigest(ref.String(), authn)
	if digest == "" {
		return "", errors.New("could not generate digest")
	}
	return fmt.Sprintf("%s@sha256:%s", ref.String(), digest), nil
}
