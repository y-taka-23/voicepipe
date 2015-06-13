package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
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
					log.Fatal(err)
				}
				if err := vp.SetupAll(); err != nil {
					log.Fatal(err)
				}
				if err := vp.BuildAll(); err != nil {
					log.Fatal(err)
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
					log.Fatal(err)
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
					log.Fatal(err)
				}
				if err := vp.CleanAll(); err != nil {
					log.Fatal(err)
				}
			},
		},
	}
	app.Run(os.Args)
}
