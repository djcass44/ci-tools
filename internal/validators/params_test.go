package validators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInternalParameterValidator(t *testing.T) {
	var cases = []struct {
		in  string
		out bool
	}{
		{
			"./testdata/valid.slsa.json",
			true,
		},
		{
			"./testdata/invalid_internal_params.slsa.json",
			false,
		},
	}
	v := &InternalParameterValidator{}
	for _, tt := range cases {
		t.Run(tt.in, func(t *testing.T) {
			ok := v.Check1(loadFile(t, tt.in))
			assert.EqualValues(t, tt.out, ok)
		})
	}
}
