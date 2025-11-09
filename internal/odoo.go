package internal

import "strings"

func DetectVersion(branch string) string {
	var version string
	if strings.HasPrefix(branch, "saas-") {
		version = strings.SplitN(branch[5:], "-", 2)[0]
		Debug.Printf("DetectVersion: detected saas branch '%v', splitting gives this '%v'", branch, "saas-"+version)
		return "saas-" + version
	}
	version = strings.SplitN(branch, "-", 2)[0]
	Debug.Printf("DetectVersion: detected regular branch '%v', splitting gives this '%v'", branch, version)
	return version
}

func isVersionBranch(branch string) bool { return branch == DetectVersion(branch) }
