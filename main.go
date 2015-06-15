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
	app.Usage = "build parameterized Docker images"
	app.Commands = []cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "fill it later",
			Action: func(c *cli.Context) {
				vp, err := NewVoicePipe(root, os.Stdout, os.Stderr)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if err := vp.SetupAll(); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if err := vp.BuildAll(); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "fill it later",
			Action: func(c *cli.Context) {
				vp, err := NewVoicePipe(root, os.Stdout, os.Stderr)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				vp.List()
			},
		},
		{
			Name:    "clean",
			Aliases: []string{"c"},
			Usage:   "fill it later",
			Action: func(c *cli.Context) {
				vp, err := NewVoicePipe(root, os.Stdout, os.Stderr)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				if err := vp.CleanAll(); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			},
		},
	}
	app.Run(os.Args)
}
