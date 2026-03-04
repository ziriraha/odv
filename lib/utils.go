package lib

import (
	"fmt"
	"maps"
	"os/exec"
	"slices"
	"sync"
)

var SortedRepoNames []string
var Repositories = make(map[string]*Repository)

func InitializeConfiguration() {
	odooHome := GetOdooPath()

	Repositories[WorkspaceRepo] = &Repository{path: odooHome + "/.vscode"}
	Repositories["community"] = &Repository{path: odooHome + "/community"}
	Repositories["enterprise"] = &Repository{path: odooHome + "/enterprise"}
	Repositories["upgrade"] = &Repository{path: odooHome + "/upgrade"}

	SortedRepoNames = slices.Sorted(maps.Keys(Repositories))
}

func PrefetchAllBranches() {
	var wg sync.WaitGroup
	for _, repo := range Repositories {
		wg.Go(func() { repo.GetBranches() })
	}
	wg.Wait()
}

func GetAllBranches() []string {
	var branches []string
	for _, repo := range Repositories {
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
