package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "voicepipe"
	app.Usage = "build parameterized Docker images"
	app.Commands = []cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "fill it later",
			Action: func(c *cli.Context) {
				BuildAction(c)
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "fill it later",
			Action: func(c *cli.Context) {
				ListAction(c)
			},
		},
		{
			Name:    "clean",
			Aliases: []string{"c"},
			Usage:   "fill it later",
			Action: func(c *cli.Context) {
				CleanAction(c)
			},
		},
	}
	app.Run(os.Args)
}
