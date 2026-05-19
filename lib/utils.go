package lib

import (
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
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
			fullPath := filepath.Join(cfg.Odoo.Path, folderName)
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

func KillProcessByPort(port int) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("invalid port number: %d", port)
	}

	pid, err := exec.Command("lsof", "-ti", fmt.Sprintf(":%d", port)).CombinedOutput()
	if err != nil || len(pid) == 0 {
		return fmt.Errorf("no process found listening on port %d", port)
	}
	pidInt, err := strconv.Atoi(strings.TrimSpace(string(pid)))
	if err != nil {
		return fmt.Errorf("invalid pid: %v", err)
	}
	process, err := os.FindProcess(pidInt)
	if err != nil {
		return fmt.Errorf("could not find process: %v", err)
	}
	err = process.Kill()
	if err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}
	return nil
}
