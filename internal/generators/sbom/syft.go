package sbom

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Snakdy/container-build-engine/pkg/oci/auth"
	"github.com/anchore/stereoscope/pkg/image"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/format/cyclonedxjson"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/google/go-containerregistry/pkg/name"
	_ "modernc.org/sqlite"
)

func Execute(ctx context.Context, bctx *civ1.BuildContext, ref, digest string) error {
	refName, err := name.ParseReference(ref)
	if err != nil {
		return fmt.Errorf("parsing reference %s: %w", ref, err)
	}
	// rewrite the reference and remove the tag so it's just
	// the digest
	ref = fmt.Sprintf("%s@sha256:%s", refName.Context(), digest)
	log.Printf("generating SBOM for ref: %s", ref)

	// configure syft to auth
	sourceOptions := syft.DefaultGetSourceConfig().WithSources("registry")
	sourceOptions.SourceProviderConfig.RegistryOptions = &image.RegistryOptions{
		Keychain: auth.KeyChain(bctx.Auth()),
	}
	var artefact *sbom.SBOM
	src, err := syft.GetSource(ctx, ref, sourceOptions)
	if err != nil {
		if !strings.Contains(err.Error(), "unknown layer media type") {
			return fmt.Errorf("fetching image source: %w", err)
		}
		log.Printf("syft failed to generate SBOM for unknown image media type: %v", err)
		// generate a skeleton SBOM for unsupported media types
		// similar to what Trivy does
		artefact = &sbom.SBOM{
			Artifacts: sbom.Artifacts{},
			Source: source.Description{
				ID:      digest,
				Name:    ref,
				Version: bctx.Repo.CommitSha,
				Metadata: map[string]any{
					"userInput":      ref,
					"tags":           bctx.Tags,
					"manifestDigest": "sha256:" + digest,
				},
			},
			Descriptor: sbom.Descriptor{
				Name:    "ci",
				Version: "unknown",
			},
		}
	}

	// hand it off to Syft
	if artefact == nil {
		artefact, err = syft.CreateSBOM(ctx, src, syft.DefaultCreateSBOMConfig())
		if err != nil {
			return fmt.Errorf("generating SBOM: %w", err)
		}
	}
	// convert the SBOM into a CycloneDX
	// report.
	encoder, err := cyclonedxjson.NewFormatEncoderWithConfig(cyclonedxjson.DefaultEncoderConfig())
	if err != nil {
		return err
	}
	var buffer bytes.Buffer
	if err := encoder.Encode(&buffer, *artefact); err != nil {
		return fmt.Errorf("encoding SBOM: %w", err)
	}
	// write to file
	if err := os.WriteFile(filepath.Join(bctx.Root, outSBOM), buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("saving SBOM: %w", err)
	}
	return nil
}
