package digestof

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile(t *testing.T) {
	hash, err := File("./testdata/file.txt")
	assert.NoError(t, err)
	// hash generated with sha256sum on linux/amd64
	assert.EqualValues(t, "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f", hash)
}
