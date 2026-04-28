package main

import (
	"os"

	"github.com/0xfakeSpike/polymarket-go/internal/cli/pmctl"
)

func main() {
	app := pmctl.App{Stdout: os.Stdout, Stderr: os.Stderr}
	os.Exit(app.Run(os.Args[1:]))
}
