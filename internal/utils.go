package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var RepositoryPaths = map[string]string{
		"community":  "/community",
		"enterprise": "/enterprise",
		"upgrade":    "/upgrade",
	}

func InitializeConfiguration() {
	odooHome := os.Getenv("ODOO_HOME")
	if len(odooHome) == 0 {
		odooHome = "."
	}
	for name := range RepositoryPaths {
		RepositoryPaths[name] = odooHome + RepositoryPaths[name]
	}
}

func DetectVersion(branch string) string {
	if strings.HasPrefix(branch, "saas-") {
		return "saas-" + strings.SplitN(branch[5:], "-", 1)[0]
	}
	return strings.SplitN(branch, "-", 1)[0]
}

func ForEachRepository(action func(repo *Repository) error) error {
	var wg sync.WaitGroup
	for _, repository := range GetRepositories() {
		wg.Add(1)
		go func(repo *Repository) {
			defer wg.Done()
			err := action(repo)
			if err != nil {
				log.Fatal(fmt.Errorf("in repository %v: %w", repo.Name, err))
			}
		}(&repository)
	}
	wg.Wait()
	return nil
}
