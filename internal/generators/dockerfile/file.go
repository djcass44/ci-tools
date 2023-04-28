package dockerfile

import (
	"io"
	"os"
)

// File copies the data from the source file
// to the destination file.
func File(src string, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	_, err = io.Copy(d, s)
	return err
}
