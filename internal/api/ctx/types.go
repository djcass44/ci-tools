package ctx

// EnvGitLabCI is used to detect that we are running inside
// a GitLab CI job.
const EnvGitLabCI = "GITLAB_CI"

// GitLabContext contains the ambient environment variables
// that are set by GitLab CI.
type GitLabContext struct {
	ProjectURL    string `env:"CI_PROJECT_URL"`
	ProjectDir    string `env:"CI_PROJECT_DIR"`
	ProjectPath   string `env:"PROJECT_PATH"`
	DefaultBranch string `env:"CI_DEFAULT_BRANCH"`

	ConfigPath string `env:"CI_CONFIG_PATH"`

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
