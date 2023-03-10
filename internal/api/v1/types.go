package v1

const (
	EnvBuildArgPrefix     = "BUILD_ARG_"
	EnvBuildDockerfile    = "BUILD_DOCKERFILE"
	EnvBuildGoImportPath  = "BUILD_GO_IMPORTPATH"
	EnvBuildImageParent   = "BUILD_IMAGE_PARENT"
	EnvBuildSLSABuildType = "BUILD_SLSA_BUILD_TYPE"

	EnvBuildCacheEnabled = "BUILD_CACHE_ENABLED"
	EnvBuildCachePath    = "BUILD_CACHE_PATH"

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
	Build map[string]BuildRecipe `json:"build"`
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
}

type BuildRepo struct {
	URL       string
	CommitSha string
	Ref       string
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
