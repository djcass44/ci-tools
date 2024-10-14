package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetImageRef(t *testing.T) {
	var cases = []struct {
		name     string
		registry string
		repo     string
		tag      string
		out      string
	}{
		{
			"normal",
			"registry.gitlab.com",
			"av1o/base-images/go-1.23-debian",
			"deadbeef",
			"registry.gitlab.com/av1o/base-images/go-1.23-debian:deadbeef",
		},
		{
			"leading slash",
			"registry.gitlab.com",
			"/av1o/base-images/go-1.23-debian",
			"deadbeef",
			"registry.gitlab.com/av1o/base-images/go-1.23-debian:deadbeef",
		},
		{
			"trailing slash",
			"registry.gitlab.com/",
			"av1o/base-images/go-1.23-debian",
			"deadbeef",
			"registry.gitlab.com/av1o/base-images/go-1.23-debian:deadbeef",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			out := getImageRef(tt.registry, tt.repo, tt.tag)
			assert.EqualValues(t, tt.out, out)
		})
	}
}
