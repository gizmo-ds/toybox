package main

import (
	"fmt"
	"os"

	"github.com/gizmo-ds/toybox/cmd/toybox/commands"
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
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}).Run(os.Args)
}
