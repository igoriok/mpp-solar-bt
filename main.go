package main

import (
	"log"
	"os"
	"watch-power-bt/cmd"
)

func init() {
	log.SetOutput(os.Stderr)
}

func main() {

	err := cmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
