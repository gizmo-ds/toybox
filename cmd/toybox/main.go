package main

import (
	"os"
	"toybox/cmd/toybox/commands"

	"github.com/urfave/cli/v2"
)

var AppVersion = "development"

func main() {
	_ = (&cli.App{
		Name:     "toybox",
		Usage:    "a toybox for learning",
		Version:  AppVersion,
		Suggest:  true,
		Commands: commands.Commands,
	}).Run(os.Args)
}
