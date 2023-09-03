package slsa

import (
	"encoding/json"
	"github.com/in-toto/in-toto-golang/in_toto"
	"github.com/in-toto/in-toto-golang/in_toto/slsa_provenance/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestVSA(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "provenance.slsa.json"))
	require.NoError(t, err)

	var provenance in_toto.ProvenanceStatementSLSA1
	require.NoError(t, json.Unmarshal(data, &provenance))

	out, err := VSA(true, &provenance, common.ProvenanceMaterial{
		URI: "https://example.com/provenances/example-1.2.3.tar.gz.intoto.jsonl",
		Digest: map[string]string{
			digestSha256: "deadbeef",
		},
	})
	assert.NoError(t, err)

	t.Log(out)
}
