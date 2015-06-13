package main

import (
	"log"
	"os"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	vp, err := NewVoicePipe(root)
	if err != nil {
		log.Fatal(err)
	}
	err = vp.Run()
	if err != nil {
		log.Fatal(err)
	}
}
