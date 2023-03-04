package v1

import (
	"bytes"
	_ "embed"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

//go:embed recipes.tpl.yaml
var defaultRecipes string

func ReadConfiguration(path string, ctx *BuildContext) (*Recipes, error) {
	var tpl *template.Template
	var err error
	if path == "" {
		log.Print("loading default configuration file")
		tpl, err = template.New("recipes.tpl.yaml").Parse(defaultRecipes)
	} else {
		log.Printf("reading configuration file: %s", path)
		tpl, err = template.New(filepath.Base(filepath.Clean(path))).ParseFiles(path)
	}
	if err != nil {
		log.Printf("failed to read configuration file: %s", err)
		return nil, err
	}
	data := new(bytes.Buffer)
	if err := tpl.Execute(data, ctx); err != nil {
		log.Printf("failed to template configuration file: %s", err)
		return nil, err
	}
	v := new(Recipes)
	if err := yaml.NewDecoder(data).Decode(&v); err != nil {
		log.Printf("failed to decode configuration file: %s", err)
		return nil, err
	}
	return v, nil
}

func WriteDockerCFG(ctx *BuildContext) error {
	path := os.Getenv("DOCKER_CONFIG")
	home, _ := os.UserHomeDir()
	if path == "" {
		path = filepath.Join(home, ".docker", "config.json")
	}
	return os.WriteFile(path, []byte(ctx.DockerCFG()), 0644)
}
