package main

import (
	"fmt"
	"github.com/djcass44/ci-tools/cmd"
	"time"
)

var (
	version = "0.0.0"
	commit  = "develop"
	date    = time.Now().String()
)

func main() {
	cmd.Execute(fmt.Sprintf("%s-%s (%s)", version, commit, date))
}
