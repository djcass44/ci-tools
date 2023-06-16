package slsa

import (
	"encoding/json"
	"fmt"
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"os"
	"path/filepath"
)

func output(ctx *v1.BuildContext, v any, digest string) error {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	// write the provenance file
	if err := os.WriteFile(filepath.Join(ctx.Root, outProvenance), data, 0644); err != nil {
		return err
	}
	// write the digest file
	if err := os.WriteFile(filepath.Join(ctx.Root, outBuild), []byte(fmt.Sprintf("%s:%s@sha256:%s", ctx.Image.Name, ctx.Repo.CommitSha, digest)), 0644); err != nil {
		return err
	}
	return nil
}
