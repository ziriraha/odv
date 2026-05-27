package lib

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"mvdan.cc/sh/v3/shell"
)

func runOdooCommand(args ...string) (string, error) {
	odooCommandString := GetConfig().Odoo.Command
	parts, err := shell.Fields(odooCommandString, nil)
	if err != nil {
		return "", err
	}
	return runCommand(parts[0], append(parts[1:], args...)...)
}

func DetectVersion(branch string) string {
	if strings.HasPrefix(branch, "saas-") {
		return "saas-" + strings.SplitN(branch[5:], "-", 2)[0]
	}
	return strings.SplitN(branch, "-", 2)[0]
}

func GetVersion(branch string) string {
	version := DetectVersion(branch)
	if strings.HasPrefix(version, "saas-") {
		return version[5:]
	}
	return version
}

func IsVersionBranch(branch string) bool {
	return branch == DetectVersion(branch)
}

func GetRemoteForBranch(branch string) string {
	if !IsVersionBranch(branch) {
		return RemoteDev
	}
	return RemoteOrigin
}

func SortBranches(branches []string) {
	slices.SortFunc(branches, func(a, b string) int {
		aVersion := GetVersion(a)
		bVersion := GetVersion(b)
		comparison := strings.Compare(aVersion, bVersion)
		if comparison != 0 {
			return -comparison
		}
		return strings.Compare(a, b)
	})
}

func (r *Repository) AutoCommit(message string) error {
	changes, err := r.GetStatus()
	if err == nil && len(changes) > 0 {
		commitMessage := fmt.Sprintf("odv auto-commit %v\n\n%s", time.Now().Format(time.RFC3339), message)
		return r.CommitAll(commitMessage)
	}
	return err
}
