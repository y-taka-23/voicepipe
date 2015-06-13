package main

import (
	"log"
)

func main() {
	vp := NewVoicePipe()
	err := vp.Run()
	if err != nil {
		log.Fatal(err)
	}
}
