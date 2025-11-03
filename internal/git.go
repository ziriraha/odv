package internal

import (
	"os/exec"
	"sync"
)

var repositories []Repository
type Repository struct {
	Name string
	Path string
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

func (r *Repository) BranchExists(branchName string) (bool) {
	cmd := exec.Command("git", "-C", r.Path, "branch", "--list", branchName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}

func (r *Repository) SwitchBranch(branchName string) error {
	cmd := exec.Command("git", "-C", r.Path, "switch", branchName)
	return cmd.Run()
}
