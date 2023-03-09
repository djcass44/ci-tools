package runtime

import (
	"fmt"
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
		newEnv = append(newEnv, fmt.Sprintf("%s=%s", k, os.ExpandEnv(v)))
	}

	// expand the arguments using a mix of the application environment
	// and our new recipe-specific environment variables
	fullEnv := append(newEnv, os.Environ()...)
	envArgs := make([]string, len(args))
	for i := range args {
		envArgs[i] = os.Expand(args[i], func(s string) string {
			return getEnv(fullEnv, s)
		})
	}
	return newEnv, envArgs
}

func getEnv(env []string, key string) string {
	for _, e := range env {
		k, v, ok := strings.Cut(e, "=")
		if !ok {
			continue
		}
		if k == key {
			return v
		}
	}
	return ""
}
