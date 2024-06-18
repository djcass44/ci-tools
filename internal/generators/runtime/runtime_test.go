package runtime

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestPrepareEnv(t *testing.T) {
	t.Run("regular env is expanded", func(t *testing.T) {
		require.NoError(t, os.Setenv("ZOO", "foo"))
		env, args := prepareEnv(map[string]string{
			"FOO": "bar",
		}, []string{
			"--foo=$FOO",
			"--zoo=$ZOO",
			"build-arg=FOO=bar",
		})
		assert.EqualValues(t, []string{"FOO=bar"}, env)
		assert.EqualValues(t, []string{"--foo=bar", "--zoo=foo", "build-arg=FOO=bar"}, args)
	})
	t.Run("env can be escaped", func(t *testing.T) {
		require.NoError(t, os.Setenv("FOO", "bar"))

		env, args := prepareEnv(map[string]string{}, []string{"-c", "export BAR=foo; echo $$BAR; echo $FOO"})
		assert.Nil(t, env)
		assert.EqualValues(t, []string{"-c", "export BAR=foo; echo $BAR; echo bar"}, args)
	})
}
