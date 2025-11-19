package internal

import (
	"io"
	"log"
	"maps"
	"os"
	"slices"
	"strings"
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

var repoNames []string
var Repositories = make(map[string]*Repository)
func InitializeConfiguration() {
	odooHome := os.Getenv("ODOO_HOME")
	if len(odooHome) == 0 {
		odooHome = "."
	}
	Debug.Printf("Configuration's Odoo Home: '%v'", odooHome)

	Repositories[".vscode"] =
			&Repository{ path: odooHome + "/.vscode", Color: color.RedString, DefaultBranch: "main" }
	Repositories["community"] =
			&Repository{ path: odooHome + "/community", Color: color.YellowString, DefaultBranch: "master" }
	Repositories["enterprise"] =
			&Repository{ path: odooHome + "/enterprise", Color: color.GreenString, DefaultBranch: "master" }
	Repositories["upgrade"] =
			&Repository{ path: odooHome + "/upgrade", Color: color.BlueString, DefaultBranch: "master" }

	repoNames = slices.Sorted(maps.Keys(Repositories))
	Debug.Printf("Initialized repositories: %v", strings.Join(repoNames, ", "))
}

func ForEachRepository(action func(i int, repoName string, repo *Repository), isConcurrent bool) {
	var wg sync.WaitGroup
	for i, n := range repoNames {
		repo := Repositories[n]
		if isConcurrent {
			wg.Add(1)
			go func(i int, n string, r *Repository) {
				defer wg.Done()
				action(i, n, r)
			}(i, n, repo)
		} else { 
			action(i, n, repo)
		}
	}
	wg.Wait()
}
