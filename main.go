package main

import (
	"log"

	"github.com/ziriraha/odoodev/cmd"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("ERROR: ")

	cmd.Execute()
}
