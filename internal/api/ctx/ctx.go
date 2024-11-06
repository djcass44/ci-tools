package ctx

import (
	"fmt"
	"github.com/Netflix/go-env"
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"os"
)

type Normalisable interface {
	Normalise() v1.BuildContext
}

func GetContext() (*v1.BuildContext, error) {
	var context Normalisable
	// check if we're in GitLab CI
	if os.Getenv(EnvGitLabCI) != "" {
		context = new(GitLabContext)
		if _, err := env.UnmarshalFromEnviron(context); err != nil {
			return nil, fmt.Errorf("parsing gitlab context from env: %w", err)
		}
	} else if os.Getenv(EnvGitHubActions) != "" {
		context = new(GitHubContext)
		if _, err := env.UnmarshalFromEnviron(context); err != nil {
			return nil, fmt.Errorf("parsing github context from env: %w", err)
		}
	}
	var buildContext v1.BuildContext
	if context != nil {
		buildContext = context.Normalise()
	}
	buildContext.Normalise()
	return &buildContext, nil
}
