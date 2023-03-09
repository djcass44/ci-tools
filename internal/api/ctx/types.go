package ctx

const EnvGitLabCI = "GITLAB_CI"

type GitLabContext struct {
	ProjectURL  string `env:"CI_PROJECT_URL"`
	ProjectDir  string `env:"CI_PROJECT_DIR"`
	ProjectPath string `env:"PROJECT_PATH"`

	JobImage     string `env:"CI_JOB_IMAGE"`
	JobID        string `env:"CI_JOB_ID"`
	JobStartedAt string `env:"CI_JOB_STARTED_AT"`

	Registry         string `env:"CI_REGISTRY"`
	RegistryImage    string `env:"CI_REGISTRY_IMAGE"`
	RegistryUser     string `env:"CI_REGISTRY_USER"`
	RegistryPassword string `env:"CI_REGISTRY_PASSWORD"`

	CommitBranch   string `env:"CI_COMMIT_BRANCH"`
	CommitTag      string `env:"CI_COMMIT_TAG"`
	CommitSha      string `env:"CI_COMMIT_SHA"`
	CommitShortSha string `env:"CI_COMMIT_SHORT_SHA"`
	CommitRefName  string `env:"CI_COMMIT_REF_NAME"`
}
