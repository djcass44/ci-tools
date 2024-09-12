package tools

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestMavenAuth(t *testing.T) {
	var cases = []struct {
		key      string
		value    string
		filename string
	}{
		{
			"mirror",
			"default=Default=https://prism.v2.dcas.dev/api/v1/maven/-/=central",
			"central.xml",
		},
		{
			"local-repo",
			"/.cache/maven",
			"local-repo.xml",
		},
		{
			"server",
			"default=default=username=password",
			"server.xml",
		},
	}

	for _, tt := range cases {
		t.Run(tt.filename, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "settings.xml")

			var localRepo string
			var servers, repos, mirrors []string
			switch tt.key {
			case "mirror":
				mirrors = append(mirrors, tt.value)
			case "server":
				servers = append(servers, tt.value)
			case "local-repo":
				localRepo = tt.value
			}

			assert.NoError(t, generate(output, localRepo, servers, repos, mirrors))
			assert.FileExists(t, output)

			src, err := os.ReadFile(filepath.Join("./testdata", tt.filename))
			require.NoError(t, err)

			dst, err := os.ReadFile(output)
			require.NoError(t, err)

			t.Logf("comparing source (%s) to destination (%s)", tt.filename, output)

			assert.EqualValues(t, string(src), string(dst))
		})
	}
}
