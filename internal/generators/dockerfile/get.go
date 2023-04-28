package dockerfile

import (
	"errors"
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
)

func Get(content *v1.DockerfileContent, path string) error {
	// prefer inline
	if content.Inline != "" {
		return Inline(content.Inline, path)
	}
	if content.File != "" {
		return File(content.File, path)
	}
	return errors.New("malformed Dockerfile descriptor - no way to retrieve it")
}
