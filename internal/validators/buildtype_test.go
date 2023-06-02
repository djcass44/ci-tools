package validators

import (
	"encoding/json"
	"github.com/djcass44/ci-tools/internal/generators/slsa"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestBuildTypeValidator_Validate(t *testing.T) {
	var cases = []struct {
		in  string
		out bool
	}{
		{
			"./testdata/valid.slsa.json",
			true,
		},
		{
			"./testdata/invalid_buildtype.slsa.json",
			false,
		},
	}
	v := BuildTypeValidator{Expected: slsa.DefaultBuildType}
	for _, tt := range cases {
		t.Run(tt.in, func(t *testing.T) {
			ok := v.Validate(loadFile(t, tt.in))
			assert.EqualValues(t, tt.out, ok)
		})
	}
}

func loadFile(t *testing.T, path string) *in_toto.ProvenanceStatementSLSA1 {
	f, err := os.Open(path)
	require.NoError(t, err)

	var statement in_toto.ProvenanceStatementSLSA1
	require.NoError(t, json.NewDecoder(f).Decode(&statement))

	return &statement
}
