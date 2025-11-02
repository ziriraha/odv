package internal

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var RepositoryPaths = map[string]string{
		"community":  "./community",
		"enterprise": "./enterprise",
		"upgrade":    "./upgrade",
	}

var (
	repositories     map[string]*git.Repository
	repositoriesOnce sync.Once
)

func GetRepositories() (map[string]*git.Repository) {
	repositoriesOnce.Do(func() {
		repositories = make(map[string]*git.Repository, len(RepositoryPaths))
		for name, path := range RepositoryPaths {
			repository, err := git.PlainOpen(path)
			if err != nil {
				log.Fatal(fmt.Errorf("opening repository %v in path %v: %v", name, path, err))
			}
			repositories[name] = repository
		}
	})
	return repositories
}

func SwitchBranch(repository *git.Repository, branchName string) error {
	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}
	return worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
	})
}