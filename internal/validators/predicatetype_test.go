package validators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPredicateTypeValidator(t *testing.T) {
	var cases = []struct {
		in  string
		out bool
	}{
		{
			"./testdata/valid.slsa.json",
			true,
		},
		{
			"./testdata/invalid_predicatetype.slsa.json",
			false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.in, func(t *testing.T) {
			ok := PredicateTypeValidator(loadFile(t, tt.in))
			assert.EqualValues(t, tt.out, ok)
		})
	}
}
