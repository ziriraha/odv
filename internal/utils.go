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
	AddRepository("community", odooHome + "/community", color.YellowString)
	AddRepository("enterprise", odooHome + "/enterprise", color.GreenString)
	AddRepository("upgrade", odooHome + "/upgrade", color.BlueString)
}

func ForEachRepository(action func(repo *Repository), isConcurrent bool) {
	var wg sync.WaitGroup
	for i := range Repositories {
		repo := &Repositories[i]
		if isConcurrent {
			wg.Add(1)
			go func(r *Repository) {
				defer wg.Done()
				action(r)
			}(repo)
		} else { action(repo) }
	}
	wg.Wait()
}
