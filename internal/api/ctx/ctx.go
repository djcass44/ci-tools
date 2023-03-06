package ctx

import (
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
	if os.Getenv("GITLAB_CI") != "" {
		context = new(GitLabContext)
		if _, err := env.UnmarshalFromEnviron(context); err != nil {
			return nil, err
		}
	}
	var buildContext v1.BuildContext
	if context != nil {
		buildContext = context.Normalise()
	}
	buildContext.Normalise()
	return &buildContext, nil
}
