package sbom

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Snakdy/container-build-engine/pkg/oci/auth"
	"github.com/anchore/stereoscope/pkg/image"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/format/cyclonedxjson"
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Execute(ctx context.Context, bctx *civ1.BuildContext, digest string) error {
	// if the image doesn't include the digest, add it
	// to the end
	ref := bctx.FQTags[0]
	if !strings.Contains(ref, "@sha256:") {
		ref = fmt.Sprintf("%s@sha256:%s", bctx.FQTags[0], digest)
	}
	log.Printf("generating SBOM for ref: %s", ref)

	// configure syft to auth
	sourceOptions := syft.DefaultGetSourceConfig()
	sourceOptions.SourceProviderConfig.RegistryOptions = &image.RegistryOptions{
		Keychain: auth.KeyChain(bctx.Auth()),
	}
	src, err := syft.GetSource(ctx, ref, sourceOptions)
	if err != nil {
		return err
	}

	// hand it off to Syft
	artefact, err := syft.CreateSBOM(ctx, src, syft.DefaultCreateSBOMConfig())
	if err != nil {
		return err
	}
	// convert the SBOM into a CycloneDX
	// report.
	encoder, err := cyclonedxjson.NewFormatEncoderWithConfig(cyclonedxjson.DefaultEncoderConfig())
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	if err := encoder.Encode(&buffer, *artefact); err != nil {
		return err
	}
	// write to file
	if err := os.WriteFile(filepath.Join(bctx.Root, outSBOM), buffer.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}
