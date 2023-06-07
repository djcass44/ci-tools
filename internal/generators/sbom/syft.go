package sbom

import (
	"github.com/anchore/stereoscope/pkg/image"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/formats"
	"github.com/anchore/syft/syft/formats/cyclonedxjson"
	"github.com/anchore/syft/syft/pkg/cataloger"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"log"
	"os"
	"path/filepath"
)

func Execute(ctx *civ1.BuildContext, digest string) error {
	// fetch the image
	ref, err := name.ParseReference(ctx.FQTags[0])
	if err != nil {
		return err
	}
	log.Printf("generating SBOM for ref: %s", ref.String())
	img, err := remote.Image(ref, remote.WithAuth(&authn.Basic{
		Username: ctx.Image.Username,
		Password: ctx.Image.Password,
	}))
	if err != nil {
		return err
	}
	syftImage := image.New(img, "", image.WithTags(ctx.Tags...))
	if err := syftImage.Read(); err != nil {
		return err
	}
	src, err := source.NewFromImage(syftImage, "")
	if err != nil {
		return err
	}

	// hand it off to Syft
	catalog, relationships, distro, err := syft.CatalogPackages(&src, cataloger.DefaultConfig())
	if err != nil {
		return err
	}
	artifact := sbom.SBOM{
		Artifacts: sbom.Artifacts{
			PackageCatalog:    catalog,
			LinuxDistribution: distro,
		},
		Relationships: relationships,
		Source: source.Metadata{
			Scheme: source.ImageScheme,
			ImageMetadata: source.ImageMetadata{
				UserInput:      ctx.Image.Name,
				Tags:           ctx.Tags,
				ManifestDigest: digest,
			},
			Name: ctx.Image.Name,
		},
		Descriptor: sbom.Descriptor{
			Name:    ctx.Image.Name,
			Version: ctx.Tags[0],
		},
	}
	// convert the SBOM into a CycloneDX
	// report.
	data, err := formats.Encode(artifact, cyclonedxjson.Format())
	if err != nil {
		return err
	}
	// write to file
	if err := os.WriteFile(filepath.Join(ctx.Root, outSBOM), data, 0644); err != nil {
		return err
	}
	return nil
}
