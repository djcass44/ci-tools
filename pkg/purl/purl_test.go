package purl

import (
	"github.com/anchore/packageurl-go"
	"github.com/stretchr/testify/assert"
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
			assert.EqualValues(t, tt.out, out)
		})
	}
}
