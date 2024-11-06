package slsa

import (
	"encoding/json"
	"fmt"
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func output(ctx *v1.BuildContext, v any, ref, digest string) error {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	// write the provenance file
	if err := os.WriteFile(filepath.Join(ctx.Root, outProvenance), data, 0644); err != nil {
		return err
	}
	// write the digest file
	imageRef := ref
	if !strings.Contains(ref, "@sha256:") {
		log.Printf("image reference '%s' contains no digest so it will be rewritten", ref)
		imageRef = fmt.Sprintf("%s@sha256:%s", ref, digest)
	}
	if err := os.WriteFile(filepath.Join(ctx.Root, outBuild), []byte(imageRef), 0644); err != nil {
		return err
	}
	return nil
}
