package v1

import (
	"bytes"
	_ "embed"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed recipes.tpl.yaml
var defaultRecipes string

func ReadConfigurations(ctx *BuildContext, paths ...string) (*Recipes, error) {
	recipes := &Recipes{
		Build:       map[string]BuildRecipe{},
		Dockerfiles: map[string]Dockerfile{},
	}
	for _, p := range paths {
		rc, err := ReadConfiguration(p, ctx)
		if err != nil {
			continue
		}
		for k, v := range rc.Build {
			recipes.Build[k] = v
		}
		for k, v := range rc.Dockerfiles {
			recipes.Dockerfiles[k] = v
		}
		log.Printf("loaded %d recipes from file: '%s'", len(rc.Build), p)
		log.Printf("loaded %d dockerfiles from file: '%s'", len(rc.Dockerfiles), p)
	}
	log.Printf("loaded %d recipes", len(recipes.Build))
	return recipes, nil
}

func ReadConfiguration(path string, ctx *BuildContext) (*Recipes, error) {
	path = strings.TrimSpace(path)
	var tpl *template.Template
	var err error
	if path == "" {
		log.Print("loading default configuration file")
		tpl, err = template.New("recipes.tpl.yaml").Parse(defaultRecipes)
	} else {
		log.Printf("reading configuration file: '%s'", path)
		tpl, err = template.New(filepath.Base(filepath.Clean(path))).ParseFiles(path)
	}
	if err != nil {
		log.Printf("failed to read configuration file: '%s'", err)
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
	path := os.Getenv(EnvDockerConfig)
	home, _ := os.UserHomeDir()
	if path == "" {
		path = filepath.Join(home, ".docker", "config.json")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(ctx.DockerCFG()), 0644)
}
