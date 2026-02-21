package lib

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
	"sync"
)

const (
	RemoteOrigin   = "origin"
	RemoteDev      = "dev"
	FallbackBranch = "master"
	WorkspaceRepo  = ".workspace"
)

type Repository struct {
	lock            sync.RWMutex
	getBranchesOnce sync.Once
	path            string
	Color           func(format string, a ...any) string
	branches        []string
}

func (r *Repository) readCommand(args ...string) (string, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	output, err := exec.Command("git", append([]string{"-C", r.path}, args...)...).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%w: %v", err, string(output))
	}
	return string(output), err
}

func (r *Repository) writeCommand(args ...string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	output, err := exec.Command("git", append([]string{"-C", r.path}, args...)...).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%w: %v", err, string(output))
	}
	return err
}

func (r *Repository) GetBranches() []string {
	if r.branches == nil {
		r.getBranchesOnce.Do(func() {
			output, err := r.readCommand("branch")
			if err == nil {
				for line := range strings.SplitSeq(output, "\n") {
					line = strings.TrimSpace(line)
					line = strings.TrimPrefix(line, "* ")
					if line != "" {
						r.branches = append(r.branches, line)
					}
				}
			}
		})
	}
	return slices.Clone(r.branches)
}

func (r *Repository) BranchExists(branchName string) bool {
	branches := r.GetBranches()
	for _, branch := range branches {
		if branch == branchName {
			return true
		}
	}
	return false
}

func (r *Repository) SwitchBranch(branchName string) error {
	return r.writeCommand("switch", branchName)
}

func (r *Repository) GetCurrentBranch() string {
	output, err := r.readCommand("branch", "--show-current")
	if err != nil {
		panic(fmt.Errorf("GetCurrentBranch error: %v", err))
	}
	return strings.TrimSpace(output)
}

func (r *Repository) GetStatus() ([]string, error) {
	var changes []string
	output, err := r.readCommand("status", "--porcelain")
	if err != nil {
		return changes, err
	}
	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimRight(line, " \t\n\r")
		if line != "" {
			changes = append(changes, line)
		}
	}
	return changes, nil
}

func (r *Repository) GetAheadBehindInfo(remote, branch string) (ahead int, behind int, err error) {
	output, err := r.readCommand("rev-list", "--left-right", "--count", fmt.Sprintf("%s...%s/%s", branch, remote, branch))
	if err != nil {
		return -1, -1, err
	}
	parts := strings.Fields(strings.TrimSpace(output))
	if len(parts) < 2 {
		return -1, -1, fmt.Errorf("unexpected rev-list output: %q", output)
	}
	fmt.Sscanf(parts[0], "%d", &ahead)
	fmt.Sscanf(parts[1], "%d", &behind)
	return ahead, behind, nil
}

func (r *Repository) Pull(remote, branch string) error {
	return r.writeCommand("pull", "--rebase", remote, branch)
}

func (r *Repository) FetchRefspec(remote, branch string) error {
	return r.writeCommand("fetch", remote, fmt.Sprintf("%s:%s", branch, branch))
}

func (r *Repository) CommitAll(message string) error {
	err := r.writeCommand("add", ".")
	if err != nil {
		return err
	}
	return r.writeCommand("commit", "-m", message)
}
