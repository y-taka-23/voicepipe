package main

import (
	"github.com/codegangsta/cli"
	"log"
)

func BuildAction(c *cli.Context, root string) {
	log.Println("BUILD")
}

func ListAction(c *cli.Context, root string) {
	log.Println("LIST")
}

func CleanAction(c *cli.Context, root string) {
	log.Println("CLEAN")
}
