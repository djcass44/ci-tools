package validators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSourceRepoValidator_Validate(t *testing.T) {
	var cases = []struct {
		in  string
		out bool
	}{
		{
			"./testdata/valid.slsa.json",
			true,
		},
		{
			"./testdata/invalid_sourcerepo.slsa.json",
			false,
		},
	}
	v := SourceRepoValidator{Expected: "pkg:github/example.org@sha1:latest"}
	for _, tt := range cases {
		t.Run(tt.in, func(t *testing.T) {
			ok := v.Validate(loadFile(t, tt.in))
			assert.EqualValues(t, tt.out, ok)
		})
	}
}
