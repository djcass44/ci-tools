package purl

import (
	"github.com/anchore/packageurl-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestParse(t *testing.T) {
	var cases = []struct {
		t      string
		in     string
		digest string
		path   string
		out    string
	}{
		{
			TypeOCI,
			"docker.io/library/busybox:latest",
			"deadbeef",
			"",
			"pkg:oci/busybox@sha256:deadbeef?repository_url=docker.io/library/busybox&tag=latest",
		},
		{
			TypeOCI,
			"docker.io/library/busybox:latest@sha256:deadbeef",
			"deadbeef",
			"",
			"pkg:oci/busybox@sha256:deadbeef?repository_url=docker.io/library/busybox&tag=latest",
		},
		{
			TypeGitLab,
			"https://gitlab.com/gitlab-org/gitlab.git",
			"deadbeef",
			"",
			"pkg:gitlab/gitlab-org/gitlab.git@sha256:deadbeef?repository_url=gitlab.com",
		},
		{
			TypeGitLab,
			"https://gitlab.com/gitlab-org/gitlab.git",
			"deadbeef",
			"lib",
			"pkg:gitlab/gitlab-org/gitlab.git@sha256:deadbeef?repository_url=gitlab.com#lib",
		},
		{
			TypeGitLab,
			"https://gitlab.com/gitlab-org/gitlab.git",
			"deadbeef",
			".",
			"pkg:gitlab/gitlab-org/gitlab.git@sha256:deadbeef?repository_url=gitlab.com",
		},
		{
			TypeGitLab,
			"https://gitlab.com/gitlab-org/gitlab.git",
			"deadbeef",
			"project",
			"pkg:gitlab/gitlab-org/gitlab.git@sha256:deadbeef?repository_url=gitlab.com#project",
		},
		{
			packageurl.TypeGeneric,
			"data.json",
			"deadbeef",
			"",
			"pkg:generic/data.json@sha256:deadbeef",
		},
	}
	for _, tt := range cases {
		t.Run(tt.in, func(t *testing.T) {
			out := Parse(tt.t, tt.in, tt.digest, "sha256", tt.path)
			actualOut, err := url.PathUnescape(out)
			require.NoError(t, err)
			assert.EqualValues(t, tt.out, actualOut)
		})
	}
}
