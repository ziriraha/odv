package main

import (
	"log"

	"github.com/ziriraha/odoodev/cmd"
	"github.com/ziriraha/odoodev/internal"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("ERROR: ")

	internal.InitializeConfiguration()

	cmd.Execute()
}
