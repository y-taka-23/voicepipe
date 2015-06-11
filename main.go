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
	}
}
