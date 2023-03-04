package ctx

import (
	v1 "github.com/djcass44/ci-tools/internal/api/v1"
	"github.com/kelseyhightower/envconfig"
	"os"
)

type Normalisable interface {
	Normalise() v1.BuildContext
}

func GetContext() (v1.BuildContext, error) {
	var context Normalisable
	// check if we're in GitLab CI
	if os.Getenv("GITLAB_CI") != "" {
		context = new(GitLabContext)
		if err := envconfig.Process("", &context); err != nil {
			return v1.BuildContext{}, err
		}
	}
	var buildContext v1.BuildContext
	if context != nil {
		buildContext = context.Normalise()
	}
	buildContext.Normalise()
	return buildContext, nil
}
