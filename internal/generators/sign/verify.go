package sign

import (
	"context"
	"errors"
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/djcass44/ci-tools/pkg/ociutil"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/sigstore/cosign/v2/cmd/cosign/cli/sign"
	"github.com/sigstore/cosign/v2/pkg/cosign"
	ociremote "github.com/sigstore/cosign/v2/pkg/oci/remote"
	"github.com/sigstore/cosign/v2/pkg/signature"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func prepare(ctx *civ1.BuildContext, target string) (name.Reference, []ociremote.Option, error) {
	ref, err := name.ParseReference(target)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("verifying image: %s", ref.String())
	var opts []ociremote.Option
	// configure authentication if the target
	// is within our registry
	if strings.HasPrefix(target, ctx.Image.Registry) {
		keychain := ociutil.KeyChain(ctx.Auth())
		opts = append(opts, ociremote.WithRemoteOptions(remote.WithAuthFromKeychain(keychain)))
	}
	ref, err = sign.GetAttachedImageRef(ref, "", opts...)
	if err != nil {
		return nil, nil, err
	}
	return ref, opts, nil
}

// VerifyAny validates that any Cosign public key in the given directory has signed
// a given image. Useful if you have multiple trusted authorities, and you want to prove
// that your image has been signed by any one of them.
func VerifyAny(ctx *civ1.BuildContext, target, dir string, offline bool) error {
	ref, opts, err := prepare(ctx, target)
	if err != nil {
		return err
	}
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) != ".pub" {
			return nil
		}
		// verify the image against this signature. Swallow any errors
		// since all but 1 should fail. If it doesn't fail (i.e. the signature is valid)
		// then return io.EOF, so we can safely return.
		if err := verify(ref, opts, path, offline); err == nil {
			log.Printf("verified signature using key: '%s'", d.Name())
			return io.EOF
		} else if err != nil {
			log.Printf("unable to verify signature using key: '%s': %s", d.Name(), err)
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}

// Verify validates that a given image has been signed by a given Cosign
// public key. It returns nil if the image is signed.
func Verify(ctx *civ1.BuildContext, target, key string, offline bool) error {
	ref, opts, err := prepare(ctx, target)
	if err != nil {
		return err
	}
	return verify(ref, opts, key, offline)
}

func verify(ref name.Reference, opts []ociremote.Option, key string, offline bool) error {
	// load the key
	verifier, err := signature.LoadPublicKey(context.Background(), key)
	if err != nil {
		return err
	}
	// fetch and verify the signatures
	log.Printf("checking if image (%s) has been signed by key: '%s'", ref.String(), key)
	signatures, _, err := cosign.VerifyImageSignatures(context.Background(), ref, &cosign.CheckOpts{
		RegistryClientOpts: opts,
		SigVerifier:        verifier,
		Offline:            offline,
		IgnoreSCT:          offline,
		IgnoreTlog:         offline,
	})
	if err != nil {
		return err
	}
	log.Printf("verified %d signature(s)", len(signatures))
	return nil
}
