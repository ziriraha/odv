package internal

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

var FallbackBranch = "master"

type Repository struct {
	lock sync.Mutex
	path string
	Color func(format string, a ...any) string
	branches []string
}

func (r *Repository) runCommand(args ...string) (string, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	output, err := exec.Command("git", append([]string{"-C", r.path}, args...)...).CombinedOutput()
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
	for _, line := range strings.Split(output, "\n") {
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
	return strings.TrimSpace(output), err
}

func (r *Repository) GetStatus() ([]string, error) {
	output, err := r.runCommand("status", "--porcelain")
	if err != nil { return nil, err }
	var changes []string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimRight(line, " \t\n\r")
		if line != "" { changes = append(changes, line) }
	}
	return changes, nil
}

func (r *Repository) GetAheadBehindInfo(branch string) (ahead int, behind int, err error) {
	remote := "dev"
	if IsVersionBranch(branch) { remote = "origin" }
	output, err := r.runCommand("rev-list", "--left-right", "--count", fmt.Sprintf("%s...%s/%s", branch, remote, branch))
	if err != nil { return 0, 0, err }
	parts := strings.Fields(strings.TrimSpace(output))
	fmt.Sscanf(parts[0], "%d", &ahead)
	fmt.Sscanf(parts[1], "%d", &behind)
	return ahead, behind, nil
}

func (r *Repository) Fetch(branch string) error {
	remote := "dev"
	if IsVersionBranch(branch) { remote = "origin" }
	_, err := r.runCommand("fetch", remote, branch)
	return err
}

func (r *Repository) Pull() error {
	_, err := r.runCommand("pull", "--ff-only")
	return err
}
