package internal

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/fatih/color"
)

var Error, Debug = log.Logger{}, log.Logger{}

func SetupLoggers(debug bool) {
	Error.SetFlags(0)
	Error.SetPrefix(color.RedString("ERROR "))
	Error.SetOutput(os.Stderr)

	Debug.SetFlags(0)
	Debug.SetPrefix(color.BlueString("DEBUG "))
	Debug.SetOutput(os.Stderr)
	if !debug { Debug.SetOutput(io.Discard) }
}

func InitializeConfiguration() {
	odooHome := os.Getenv("ODOO_HOME")
	if len(odooHome) == 0 {
		odooHome = "."
	}
	Debug.Printf("Configuration's Odoo Home: '%v'", odooHome)
	AddRepository(".vscode", odooHome + "/.vscode", color.RedString, "main")
	AddRepository("community", odooHome + "/community", color.YellowString, "master")
	AddRepository("enterprise", odooHome + "/enterprise", color.GreenString, "master")
	AddRepository("upgrade", odooHome + "/upgrade", color.BlueString, "master")
}

func ForEachRepository(action func(i int,repo *Repository), isConcurrent bool) {
	var wg sync.WaitGroup
	for i := range Repositories {
		repo := &Repositories[i]
		if isConcurrent {
			wg.Add(1)
			go func(i int, r *Repository) {
				defer wg.Done()
				action(i, r)
			}(i, repo)
		} else { action(i, repo) }
	}
	wg.Wait()
}
