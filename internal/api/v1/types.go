package v1

type BuildRecipe struct {
	DockerCFG bool              `yaml:"dockercfg,omitempty"`
	Env       map[string]string `yaml:"env,omitempty"`
	Args      []string          `yaml:"args,omitempty"`
}

type Recipes struct {
	Build map[string]BuildRecipe `json:"build"`
}

type BuildContext struct {
	Root       string
	Context    string
	Image      ImageConfig
	Tags       []string
	FQTags     []string
	Dockerfile DockerfileConfig
}

type ImageConfig struct {
	Parent   string
	Name     string
	Registry string
	Username string
	Password string
}

type DockerfileConfig struct {
	File string
	Args []string
}
