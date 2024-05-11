package main

import (
	"fmt"
	"os"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
)

func mutateHook(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	for _, model := range b.Models {
		for _, field := range model.Fields {
			name := field.Name
			if name == "id" {
				name = "_id,omitempty"
			}
			// add bson tag to the field
			field.Tag += `bson:"` + name + `"`
		}
	}
	return b
}

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load config", err.Error())
		os.Exit(2)
	}

	p := modelgen.Plugin{
		MutateHook: mutateHook,
	}

	err = api.Generate(
		cfg,
		api.ReplacePlugin(&p),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}
