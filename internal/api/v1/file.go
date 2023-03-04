package v1

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"path/filepath"
)

func ReadConfiguration(path string, ctx *BuildContext) (*Recipes, error) {
	log.Printf("reading configuration file: %s", path)
	tpl, err := template.New(filepath.Base(filepath.Clean(path))).ParseFiles(path)
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
