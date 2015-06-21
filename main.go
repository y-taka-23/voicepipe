package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	app := cli.NewApp()
	app.Name = "voicepipe"
	app.Usage = "Build parameterized Docker images from a single Dockerfile"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "Builds parameterized Docker images",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "latest, L",
					Usage: "build an original Dockerfile as the latest image",
				},
			},
			Action: func(c *cli.Context) {
				vp, err := newVoicePipe(root, os.Stdout, os.Stderr)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if err := vp.setupAll(); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if err := vp.buildAll(); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if c.Bool("latest") {
					if err := vp.buildLatest(); err != nil {
						fmt.Fprintln(os.Stderr, err)
						os.Exit(1)
					}
				}
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "Shows a list of tags",
			Action: func(c *cli.Context) {
				vp, err := newVoicePipe(root, os.Stdout, os.Stderr)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				vp.list()
			},
		},
		{
			Name:    "clean",
			Aliases: []string{"c"},
			Usage:   "Cleans temporary files up",
			Action: func(c *cli.Context) {
				vp, err := newVoicePipe(root, os.Stdout, os.Stderr)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if err := vp.cleanAll(); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			},
		},
	}
	app.Run(os.Args)
}
