package lib

import (
	"os"
	"slices"
	"strings"
)

func GetOdooPath() string {
	odooHome := os.Getenv("ODOO_HOME")
	if len(odooHome) == 0 {
		odooHome = "."
	}
	return odooHome
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
