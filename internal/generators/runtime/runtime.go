package runtime

import (
	"fmt"
	"github.com/a8m/envsubst"
	"github.com/a8m/envsubst/parse"
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Execute(ctx *civ1.BuildContext, r *civ1.BuildRecipe) error {
	env, envArgs := prepareEnv(r.Env, r.Args)

	// run the command
	log.Printf("running command: [%s %s] with env: [%s]", r.Command, strings.Join(envArgs, " "), strings.Join(env, " "))
	cmd := exec.Command(r.Command, envArgs...) //nolint:gosec
	// only change directories if
	// the command requires it
	if r.CD {
		cmd.Dir = filepath.Join(ctx.Root, ctx.Context)
	} else {
		cmd.Dir = ctx.Root
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Environ(), env...)

	if err := cmd.Run(); err != nil {
		log.Printf("command execution failed: %s", err)
		return err
	}
	return nil
}

func prepareEnv(env map[string]string, args []string) ([]string, []string) {
	// expand our new environment variables
	var newEnv []string
	for k, v := range env {
		val, _ := envsubst.String(v)
		newEnv = append(newEnv, fmt.Sprintf("%s=%s", k, val))
	}

	// expand the arguments using a mix of the application environment
	// and our new recipe-specific environment variables
	fullEnv := append(newEnv, os.Environ()...)
	parser := parse.New("", fullEnv, parse.Relaxed)

	envArgs := make([]string, len(args))
	for i := range args {
		val, _ := parser.Parse(args[i])
		envArgs[i] = val
	}
	return newEnv, envArgs
}
