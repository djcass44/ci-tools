package runtime

import (
	"fmt"
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Execute(r *v1.BuildRecipe) error {
	envArgs := make([]string, len(r.Args))
	for i := range r.Args {
		envArgs[i] = os.ExpandEnv(r.Args[i])
	}

	var env []string
	for k, v := range r.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, os.ExpandEnv(v)))
	}

	// run the command
	log.Printf("running command: [%s %s] with env: [%s]", r.Command, strings.Join(envArgs, " "), strings.Join(env, " "))
	cmd := exec.Command(r.Command, envArgs...) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Environ(), env...)

	if err := cmd.Run(); err != nil {
		log.Printf("command execution failed: %s", err)
		return err
	}
	return nil
}
