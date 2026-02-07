package lib

import (
	"io"
	"log"
	"maps"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/fatih/color"
)

func GrayString(format string, a ...any) string {
	return color.New(color.FgBlack, color.FgWhite).Sprintf(format, a...)
}

var Error, Debug = log.Logger{}, log.Logger{}
func SetupLoggers(debug bool) {
	Error.SetFlags(0)
	Error.SetPrefix(color.RedString("ERROR "))
	Error.SetOutput(os.Stderr)

	Debug.SetFlags(0)
	Debug.SetPrefix(color.BlueString("DEBUG "))
	Debug.SetOutput(os.Stderr)
	if !debug { Debug.SetOutput(io.Discard) }
}

var repoNames []string
var Repositories = make(map[string]*Repository)
func InitializeConfiguration() {
	odooHome := GetOdooPath()
	Debug.Printf("Odoo Home: '%v'", odooHome)

	Repositories[".workspace"] = &Repository{ path: odooHome + "/.vscode", Color: color.RedString }
	Repositories["community"] = &Repository{ path: odooHome + "/community", Color: color.YellowString }
	Repositories["enterprise"] = &Repository{ path: odooHome + "/enterprise", Color: color.GreenString }
	Repositories["upgrade"] = &Repository{ path: odooHome + "/upgrade", Color: color.BlueString }

	repoNames = slices.Sorted(maps.Keys(Repositories))
	Debug.Printf("Initialized repositories: %v", strings.Join(repoNames, ", "))
}

func ForEachRepository(action func(i int, repoName string, repo *Repository) error, isConcurrent bool) map[string]error {
	var wg sync.WaitGroup
	var err sync.Map
	for i, n := range repoNames {
		repo := Repositories[n]
		if isConcurrent { wg.Go(func() { err.Store(n, action(i, n, repo)) })
		} else { err.Store(n, action(i, n, repo)) }
	}
	wg.Wait()
	var errors = make(map[string]error)
	err.Range(func(key, value any) bool {
		if value != nil { errors[key.(string)] = value.(error) }
		return true
	})
	return errors
}

func PrintRepositoryErrors(errors map[string]error) {
	for repoName, err := range errors {
		Error.Printf("in repository %v: %v", Repositories[repoName].Color(repoName), err)
	}
}
