package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
)

func InitializeConfiguration() {
	odooHome := os.Getenv("ODOO_HOME")
	if len(odooHome) == 0 {
		odooHome = "."
	}
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

func ForEachRepository(action func(repo *Repository) error) error {
	var wg sync.WaitGroup
	for i := range Repositories {
		wg.Add(1)
		repo := &Repositories[i]
		go func(r *Repository) {
			defer wg.Done()
			err := action(r)
			if err != nil {
				log.Fatal(fmt.Errorf("in repository %v: %w", r.Color(r.Name), err))
			}
		}(repo)
	}
	wg.Wait()
	return nil
}
