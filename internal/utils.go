package internal

import (
	"os"
	"strings"
)

var RepositoryPaths = map[string]string{
		"community":  "/community",
		"enterprise": "/enterprise",
		"upgrade":    "/upgrade",
	}

func InitializeConfiguration() {
	odooHome := os.Getenv("ODOO_HOME")
	if len(odooHome) == 0 {
		odooHome = "."
	}
	for name := range RepositoryPaths {
		RepositoryPaths[name] = odooHome + RepositoryPaths[name]
	}
}

func DetectVersion(branch string) string {
	if strings.HasPrefix(branch, "saas-") {
		return "saas-" + strings.SplitN(branch[5:], "-", 1)[0]
	}
	return strings.SplitN(branch, "-", 1)[0]
}
