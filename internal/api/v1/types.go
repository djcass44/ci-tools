package v1

const (
	EnvBuildExtraArgs = "BUILD_EXTRA_ARGS"
	EnvBuildArgPrefix = "BUILD_ARG_"
	// EnvBuildTags is a comma-separated list of
	// tags that should be added to the image. Tags
	// are only added on the trunk.
	EnvBuildTags = "BUILD_TAGS"
	// EnvBuildDockerfile dictates where Dockerfile-based recipes
	// should look for the Dockerfile
	EnvBuildDockerfile = "BUILD_DOCKERFILE"
	// EnvBuildGoImportPath dictates where Go-based recipes should
	// look for the main.go file
	EnvBuildGoImportPath = "BUILD_GO_IMPORTPATH"
	// EnvBuildImageParent dictates what the build tool should use as a base layer.
	// Some recipes (e.g. BuildKit) will ignore this.
	EnvBuildImageParent = "BUILD_IMAGE_PARENT"
	// EnvBuildSLSABuildType overrides the SLSA build type field.
	EnvBuildSLSABuildType = "BUILD_SLSA_BUILD_TYPE"

	// EnvBuildCacheEnabled dictates whether we should instruct build tools
	// to cache files. Disabling this may fix some problems at the cost of performance.
	EnvBuildCacheEnabled = "BUILD_CACHE_ENABLED"
	// EnvBuildCachePath describes the directory that
	// we should instruct build tools to store their temporary files.
	EnvBuildCachePath = "BUILD_CACHE_PATH"

	EnvDockerConfig = "DOCKER_CONFIG"
)

const (
	DefaultDockerfile = "Dockerfile"
	DefaultCacheName  = ".cache"
)

type BuildRecipe struct {
	DockerCFG bool              `yaml:"dockercfg,omitempty"`
	CD        bool              `yaml:"cd,omitempty"`
	Env       map[string]string `yaml:"env,omitempty"`
	Command   string            `yaml:"command"`
	Args      []string          `yaml:"args,omitempty"`
}

type Recipes struct {
	Build       map[string]BuildRecipe `yaml:"build"`
	Dockerfiles map[string]Dockerfile  `yaml:"dockerfiles"`
}

type BuildContext struct {
	Builder    string
	BuildID    string
	Root       string
	Context    string
	Image      ImageConfig
	Tags       []string
	FQTags     []string
	Dockerfile DockerfileConfig
	Go         GoConfig
	Repo       BuildRepo
	Cache      BuildCache
	StartTime  string
	Provider   string
	ConfigPath string
	ExtraArgs  []string
}

type BuildRepo struct {
	URL       string
	CommitSha string
	Ref       string
	Trunk     bool
}

type BuildCache struct {
	Enabled bool
	Path    string
}

type ImageConfig struct {
	Parent   string
	Base     string
	Name     string
	Registry string
	Username string
	Password string
}

type DockerfileConfig struct {
	File string
	Args []string
}

type GoConfig struct {
	ImportPath string
}

type Dockerfile struct {
	Content DockerfileContent `yaml:"content"`
}

type DockerfileContent struct {
	Inline string `yaml:"inline,omitempty"`
	File   string `yaml:"file,omitempty"`
}
