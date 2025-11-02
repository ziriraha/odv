package internal

import (
	"os/exec"
)

var RepositoryPaths = map[string]string{
		"community":  "./community",
		"enterprise": "./enterprise",
		"upgrade":    "./upgrade",
	}

type Repository struct {
	Name string
	Path string
}

func GetRepositories() [3]Repository {
	var repositories [3]Repository
	i := 0
	for name, path := range RepositoryPaths {
		repositories[i] = Repository{
			Name: name,
			Path: path,
		}
		i++
	}
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