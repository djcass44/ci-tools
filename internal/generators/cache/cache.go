package cache

import (
	civ1 "github.com/djcass44/ci-tools/internal/api/v1"
	"log"
	"os"
)

func Execute(ctx *civ1.BuildContext) {
	if !ctx.Cache.Enabled {
		log.Printf("skipping cache creation")
		return
	}
	log.Printf("creating cache directory: %s", ctx.Cache.Path)
	// we need to grant all permissions here because if the container
	// image uses a different uid then we will run into issues here
	if err := os.MkdirAll(ctx.Cache.Path, 0777); err != nil {
		log.Printf("failed to create cache directory: %s", err)
	}
}
