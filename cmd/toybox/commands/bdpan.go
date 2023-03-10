package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gizmo-ds/toybox/internal/utils"
	"github.com/gizmo-ds/toybox/pkg/bdpan"
	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

func init() {
	Commands = append(Commands, BDPanTools)
}

type bdpanConfig struct {
	Bduss string `json:"bduss"`
	Name  string `json:"name"`
}

var BDPanTools = &cli.Command{
	Name:  "bdpan",
	Usage: "baidu pan tools",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			DefaultText: "<UserConfigDir>/toybox/bdpan.json",
		},
	},
	Before: func(c *cli.Context) error {
		configFile := c.String("config")
		if configFile == "" {
			configDir := utils.ConfigDir()
			configFile = filepath.Join(configDir, "bdpan.json")
			_ = c.Set("config", configFile)
		}
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			return os.WriteFile(configFile, []byte("{}"), 0644)
		}
		return nil
	},
	Subcommands: []*cli.Command{
		{
			Name:  "whoami",
			Usage: "currently logged in user",
			Action: func(c *cli.Context) error {
				configFile := c.String("config")
				data, err := os.ReadFile(configFile)
				if err != nil {
					return err
				}
				var config bdpanConfig
				if err = json.Unmarshal(data, &config); err != nil {
					return err
				}
				if config.Bduss == "" || config.Name == "" {
					return errors.New("not logged in")
				}
				fmt.Println(color.FgGreen.Render("Logged in as"), color.FgYellow.Render(config.Name))
				return nil
			},
		},
		{
			Name:  "login",
			Usage: "login to baidu pan",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "bduss",
					Usage: "BDUSS cookie value",
				},
				&cli.StringFlag{
					Name:  "cookies",
					Usage: "full cookies string",
				},
			},
			Action: func(c *cli.Context) error {
				bduss := c.String("bduss")
				cookies := c.String("cookies")
				if c.IsSet("cookies") && bduss == "" {
					r, _ := regexp.Compile(`BDUSS=([^;]+)`)
					bduss = r.FindStringSubmatch(cookies)[1]
				} else if bduss == "" {
					bduss = c.Args().First()
				}
				if bduss == "" {
					return errors.New("BDUSS is empty")
				}

				name, err := bdpan.CheckBduss(bduss)
				if err != nil {
					return err
				}
				fmt.Println(color.FgGreen.Render("Login success!"), "Welcome", color.FgYellow.Render(name))

				config := bdpanConfig{Bduss: bduss, Name: name}
				data, err := json.MarshalIndent(config, "", "  ")
				if err != nil {
					return err
				}
				return os.WriteFile(c.String("config"), data, 0644)
			},
		},
		{
			// 获取文件的rapid
			Name:    "rapid",
			Usage:   "get rapid of file",
			Aliases: []string{"id"},
			Action: func(c *cli.Context) error {
				filename := c.Args().First()
				if filename == "" {
					return errors.New("filename is empty")
				}
				file, err := os.Open(filename)
				if err != nil {
					return err
				}
				name, md5, size, err := bdpan.GetFileHash(file)
				if err != nil {
					return err
				}
				fmt.Print(strings.Join([]string{md5, strconv.Itoa(size), name}, "#"))
				return nil
			},
		},
		{
			Name:    "rapidupload",
			Usage:   "rapid upload file to baidu pan",
			Aliases: []string{"ru"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "md5",
					Aliases: []string{"m"},
					Usage:   "file md5",
				},
				&cli.IntFlag{
					Name:        "size",
					Aliases:     []string{"z"},
					Usage:       "file size in bytes",
					DefaultText: "auto detect from rapid",
				},
				&cli.StringFlag{
					Name:    "rapid",
					Aliases: []string{"r"},
				},
			},
			Action: func(c *cli.Context) error {
				configFile := c.String("config")
				data, err := os.ReadFile(configFile)
				if err != nil {
					return err
				}
				var config bdpanConfig
				if err = json.Unmarshal(data, &config); err != nil {
					return err
				}
				if config.Bduss == "" || config.Name == "" {
					return errors.New("not logged in")
				}

				savename := c.Args().First()
				md5 := c.String("md5")
				size := c.Int("size")
				rapid := c.String("rapid")
				if md5 == "" || size <= 0 {
					if rapid != "" {
						arr := strings.Split(rapid, "#")
						if len(arr) < 3 {
							return errors.New("invalid rapid")
						}
						md5 = arr[0]
						var err error
						if len(arr[1]) == 32 {
							size, err = strconv.Atoi(arr[2])
							if err != nil {
								return errors.New("invalid rapid")
							}
						} else {
							size, err = strconv.Atoi(arr[1])
							if err != nil {
								return errors.New("invalid rapid")
							}
						}
					}
				}
				if savename == "" {
					if rapid != "" {
						arr := strings.Split(rapid, "#")
						savename = "/" + arr[len(arr)-1]
					} else {
						savename = "/toybox_ru_" + time.Now().Format("20060102150405")
					}
				} else if strings.HasSuffix(savename, "/") && rapid != "" {
					arr := strings.Split(rapid, "#")
					savename += arr[len(arr)-1]
				}
				if err = bdpan.RapidUpload(config.Bduss, savename, md5, size, 3); err != nil {
					return err
				}
				fmt.Println(color.FgGreen.Render("Rapid upload success!"), "Saved as", color.FgYellow.Render(savename))
				return nil
			},
		},
	},
}
