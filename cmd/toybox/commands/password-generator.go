package commands

import (
	"fmt"

	pwdgen "github.com/gizmo-ds/toybox/pkg/pwd-gen"
	"github.com/urfave/cli/v2"
)

func init() {
	Commands = append(Commands, PasswordGenerator)
}

var PasswordGenerator = &cli.Command{
	Name:  "pwd",
	Usage: "generate password",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "length",
			Aliases: []string{"l"},
			Usage:   "password length",
			Value:   16,
		},
		&cli.BoolFlag{
			Name:    "uppercase",
			Usage:   "include uppercase letters",
			Aliases: []string{"up"},
			Value:   true,
		},
		&cli.BoolFlag{
			Name:    "bighead",
			Usage:   "first letter is uppercase",
			Aliases: []string{"bh"},
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "lowercase",
			Aliases: []string{"low"},
			Usage:   "include lowercase letters",
			Value:   true,
		},
		&cli.BoolFlag{
			Name:    "ambiguous",
			Aliases: []string{"a"},
			Usage:   "include ambiguous characters",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "numbers",
			Aliases: []string{"n"},
			Usage:   "include numbers",
			Value:   true,
		},
		&cli.IntFlag{
			Name:    "min-number",
			Aliases: []string{"mn"},
			Usage:   "minimum number of numbers",
			Value:   1,
		},
		&cli.BoolFlag{
			Name:    "special",
			Aliases: []string{"s"},
			Usage:   "include special characters",
			Value:   true,
		},
		&cli.IntFlag{
			Name:    "min-special",
			Aliases: []string{"ms"},
			Usage:   "minimum number of special characters",
			Value:   1,
		},
		&cli.BoolFlag{
			Name:  "color",
			Usage: "colorize output",
			Value: true,
		},
	},
	Action: func(c *cli.Context) error {
		p, pc := pwdgen.Generate(pwdgen.Option{
			Length:     c.Int("length"),
			Ambiguous:  c.Bool("ambiguous"),
			Uppercase:  c.Bool("uppercase"),
			Lowercase:  c.Bool("lowercase"),
			Numbers:    c.Bool("numbers"),
			MinNumber:  c.Int("min-number"),
			Special:    c.Bool("special"),
			MinSpecial: c.Int("min-special"),
			BigHead:    c.Bool("bighead"),
		})
		if c.Bool("color") {
			fmt.Println(pc)
		} else {
			fmt.Println(p)
		}
		return nil
	},
}
