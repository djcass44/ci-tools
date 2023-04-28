package dockerfile

import (
	"os"
)

// Inline writes the value of the source
// to the destination file.
func Inline(src string, dst string) error {
	return os.WriteFile(dst, []byte(src), 0644)
}
