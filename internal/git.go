package internal

import (
	"os/exec"
	"strings"
	"sync"
)

var repositories []Repository
type Repository struct {
	Name string
	Path string
	branches []string
}

var repositoriesOnce sync.Once
func GetRepositories() []Repository {
	repositoriesOnce.Do(func() {
		for name, path := range RepositoryPaths {
			repositories = append(repositories, Repository{
				Name: name,
				Path: path,
			})
		}
	})
	return repositories
}

func (r *Repository) GetBranches() ([]string, error) {
	if r.branches != nil {
		return r.branches, nil
	}
	output, err := exec.Command("git", "-C", r.Path, "branch").Output()
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "* ")
		if line != "" {
			r.branches = append(r.branches, line)
		}
	}
	return r.branches, nil
}

func (r *Repository) BranchExists(branchName string) (bool) {
	output, err := exec.Command("git", "-C", r.Path, "branch", "--list", branchName).Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}

func (r *Repository) SwitchBranch(branchName string) error {
	return exec.Command("git", "-C", r.Path, "switch", branchName).Run()
}

func (r *Repository) GetCurrentBranch() (string, error) {
	output, err := exec.Command("git", "-C", r.Path, "branch", "--show-current").Output()
	return strings.TrimSpace(string(output)), err
}

func (r *Repository) Fetch(branch string) error {
	remote := "dev"
	if isVersionBranch(branch) {
		remote = "origin"
	}
	return exec.Command("git", "-C", r.Path, "fetch", remote, branch).Run()
}

func (r *Repository) Pull() error {
	return exec.Command("git", "-C", r.Path, "pull", "--ff-only").Run()
}
