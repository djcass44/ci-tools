package ociutil

import (
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
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
func GetDigest(target string) string {
	sha, err := crane.Digest(target, crane.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		log.Printf("unable to generate digest for: %s", target)
		return ""
	}
	return strings.TrimPrefix(sha, "sha256:")
}

func SnapshotImage(target string) (string, error) {
	ref, err := name.ParseReference(target)
	if err != nil {
		return "", err
	}
	if d, err := name.NewDigest(ref.String()); err == nil {
		return d.String(), nil
	}
	if _, err := name.NewTag(ref.String()); err != nil {
		return "", err
	}
	digest := GetDigest(ref.String())
	return fmt.Sprintf("%s@sha256:%s", ref.String(), digest), nil
}
