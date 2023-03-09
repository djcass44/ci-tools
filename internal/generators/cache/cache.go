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
	if err := os.MkdirAll(ctx.Cache.Path, 0750); err != nil {
		log.Printf("failed to create cache directory: %s", err)
	}
}
