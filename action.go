package main

import (
	"github.com/codegangsta/cli"
	"log"
)

func BuildAction(c *cli.Context) {
	log.Println("BUILD")
}

func ListAction(c *cli.Context) {
	log.Println("LIST")
}

func CleanAction(c *cli.Context) {
	log.Println("CLEAN")
}
