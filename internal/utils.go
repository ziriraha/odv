package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var Error = log.Logger{}
var Debug = log.Logger{}

func SetupLoggers(debug bool) {
	Error.SetFlags(0)
	Error.SetPrefix(color.RedString("ERROR "))
	Error.SetOutput(os.Stderr)

	Debug.SetFlags(0)
	Debug.SetPrefix(color.BlueString("DEBUG "))
	Debug.SetOutput(os.Stderr)
	if !debug {
		Debug.SetOutput(io.Discard)
	}
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

func DetectVersion(branch string) string {
	if strings.HasPrefix(branch, "saas-") {
		return "saas-" + strings.SplitN(branch[5:], "-", 1)[0]
	}
	return strings.SplitN(branch, "-", 1)[0]
}

func isVersionBranch(branch string) bool {
	return branch == DetectVersion(branch)
}

func ForEachRepository(action func(repo *Repository) error, isConcurrent bool) error {
	var wg sync.WaitGroup
	for i := range Repositories {
		repo := &Repositories[i]
		if isConcurrent {
			wg.Add(1)
			go func(r *Repository) {
				defer wg.Done()
				err := action(r)
				if err != nil {
					Error.Fatal(fmt.Errorf("in repository %v: %w", r.Color(r.Name), err))
				}
			}(repo)
		} else {
			err := action(repo)
			if err != nil {
				Error.Fatal(fmt.Errorf("in repository %v: %w", repo.Color(repo.Name), err))
			}
		}
	}
	wg.Wait()
	return nil
}
