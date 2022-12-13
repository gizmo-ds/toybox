package main

import (
	"os"
	"toybox/cmd/toybox/commands"

	"github.com/urfave/cli/v2"
)

var AppVersion = "v0.1.0"

func main() {
	_ = (&cli.App{
		Name:    "toybox",
		Usage:   "a toybox for learning",
		Version: AppVersion,
		Suggest: true,
		Commands: []*cli.Command{
			commands.PasswordGenerator,
		},
	}).Run(os.Args)
}
