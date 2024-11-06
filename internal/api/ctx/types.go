package ctx

// EnvGitLabCI is used to detect that we are running inside
// a GitLab CI job.
const EnvGitLabCI = "GITLAB_CI"

// EnvGitHubActions is used to detect that we are running
// inside a GitHub Action
const EnvGitHubActions = "GITHUB_ACTIONS"

// GitLabContext contains the ambient environment variables
// that are set by GitLab CI.
type GitLabContext struct {
	ProjectURL          string `env:"CI_PROJECT_URL"`
	ProjectDir          string `env:"CI_PROJECT_DIR"`
	ProjectPath         string `env:"PROJECT_PATH"`
	ProjectPathOverride string `env:"PROJECT_PATH_OVERRIDE"`
	DefaultBranch       string `env:"CI_DEFAULT_BRANCH"`

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

type GitHubContext struct {
	ServerURL string `env:"GITHUB_SERVER_URL"`

	ProjectDir          string `env:"GITHUB_WORKSPACE"`
	ProjectNamespace    string `env:"GITHUB_REPOSITORY"`
	ProjectPath         string `env:"PROJECT_PATH"`
	ProjectPathOverride string `env:"PROJECT_PATH_OVERRIDE"`

	JobID string `env:"GITHUB_RUN_ID"`

	RegistryUsername string `env:"GITHUB_ACTOR"`
	RegistryPassword string `env:"GITHUB_TOKEN"`

	CommitSha     string `env:"GITHUB_SHA"`
	CommitRefName string `env:"GITHUB_REF_NAME"`
	CommitRefType string `env:"GITHUB_REF_TYPE"`
}
