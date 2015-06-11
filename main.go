package main

import (
	"log"
	"os"
)

func main() {
	vp := NewVoicePipe()
	err := vp.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	os.Exit(0)
}
