package internal

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type Repository struct {
	lock sync.Mutex

	path string
	Color func(format string, a ...any) string

	branches []string
	DefaultBranch string
}

func (r *Repository) runCommand(args ...string) (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	output, err := exec.Command("git", append([]string{"-C", r.path}, args...)...).Output()
	if err != nil { err = fmt.Errorf("%w: %v", err, string(output)) }
	return string(output), err
}

func (r *Repository) GetBranches() ([]string) {
	if r.branches != nil { return r.branches }
	output, err := r.runCommand("branch")
	if err != nil {
		Debug.Printf("GetBranches error: %v", err)
		return nil
	}
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "* ")
		if line != "" { r.branches = append(r.branches, line) }
	}
	return r.branches
}

func (r *Repository) BranchExists(branchName string) (bool) {
	output, err := r.runCommand("branch", "--list", branchName)
	if err != nil {
		Debug.Printf("BranchExists error: %v", err)
		return false
	}
	return len(output) > 0
}

func (r *Repository) SwitchBranch(branchName string) error {
	_, err := r.runCommand("switch", branchName)
	return err
}

func (r *Repository) GetCurrentBranch() (string, error) {
	output, err := r.runCommand("branch", "--show-current")
	return strings.TrimSpace(string(output)), err
}

func (r *Repository) Fetch(branch string) error {
	remote := "dev"
	if isVersionBranch(branch) { remote = "origin" }
	_, err := r.runCommand("fetch", remote, branch)
	return err
}

func (r *Repository) Pull() error {
	_, err := r.runCommand("pull", "--ff-only")
	return err
}
