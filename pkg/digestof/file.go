package digestof

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("failed to open file: %s: %s", path, err)
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Printf("failed to generate digest of file: %s: %s", path, err)
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
