package runtime

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPrepareEnv(t *testing.T) {
	_ = os.Setenv("ZOO", "foo")
	env, args := prepareEnv(map[string]string{
		"FOO": "bar",
	}, []string{
		"--foo=$FOO",
		"--zoo=$ZOO",
		"build-arg=FOO=bar",
	})
	assert.EqualValues(t, []string{"FOO=bar"}, env)
	assert.EqualValues(t, []string{"--foo=bar", "--zoo=foo", "build-arg=FOO=bar"}, args)
}
