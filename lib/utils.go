package lib

import (
	"fmt"
	"maps"
	"os/exec"
	"path/filepath"
	"slices"
	"sync"
)

var (
	sortedRepoNames     []string
	sortedRepoNamesOnce sync.Once

	repositories     = make(map[string]*Repository)
	repositoriesOnce sync.Once
)

func GetRepositories() map[string]*Repository {
	repositoriesOnce.Do(func() {
		cfg := GetConfig()

		for name, folderName := range cfg.Repositories {
			fullPath := filepath.Join(cfg.OdooHome, folderName)
			repositories[name] = &Repository{path: fullPath}
		}

		var wg sync.WaitGroup
		for _, repo := range repositories {
			wg.Go(func() { repo.GetBranches() })
		}
		wg.Wait()
	})
	return repositories
}

func GetRepository(name string) *Repository {
	repo, exists := GetRepositories()[name]
	if !exists {
		panic("Repository not found: " + name)
	}
	return repo
}

func GetSortedRepoNames() []string {
	sortedRepoNamesOnce.Do(func() {
		sortedRepoNames = slices.Sorted(maps.Keys(GetRepositories()))
	})
	return sortedRepoNames
}

func GetRepositoryByIndex(index int) *Repository {
	repoNames := GetSortedRepoNames()
	if index < 0 || index >= len(repoNames) {
		panic(fmt.Sprintf("Repository index out of range: %d", index))
	}
	return GetRepository(repoNames[index])
}

func GetAllBranches() []string {
	var branches []string
	for repoName, repo := range GetRepositories() {
		if repoName == ".workspace" {
			continue // skip as .workspace will create branches for everything.
		}
		for _, branch := range repo.GetBranches() {
			branches = append(branches, branch)
		}
	}
	SortBranches(branches)
	return slices.Compact(branches)
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%w: %v", err, string(output))
	}
	return string(output), err
}
