package sign

import (
	"context"
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/sigstore/cosign/cmd/cosign/cli/sign"
	"github.com/sigstore/cosign/pkg/cosign"
	ociremote "github.com/sigstore/cosign/pkg/oci/remote"
	"github.com/sigstore/cosign/pkg/signature"
	"log"
	"strings"
)

func Verify(ctx *civ1.BuildContext, target, key string) error {
	ref, err := name.ParseReference(target)
	if err != nil {
		return err
	}
	log.Printf("verifying image: %s", ref.String())
	var opts []ociremote.Option
	// configure authentication if the target
	// is within our registry
	if strings.HasPrefix(target, ctx.Image.Registry) {
		opts = append(opts, ociremote.WithRemoteOptions(remote.WithAuth(&authn.Basic{
			Username: ctx.Image.Username,
			Password: ctx.Image.Password,
		})))
	}
	ref, err = sign.GetAttachedImageRef(ref, "", opts...)
	if err != nil {
		return err
	}
	// load the key
	verifier, err := signature.LoadPublicKey(context.TODO(), key)
	if err != nil {
		return err
	}
	// fetch and verify the signatures
	signatures, _, err := cosign.VerifyImageSignatures(context.TODO(), ref, &cosign.CheckOpts{
		RegistryClientOpts: opts,
		SigVerifier:        verifier,
	})
	if err != nil {
		return err
	}
	log.Printf("verified %d signature(s)", len(signatures))
	return nil
}
