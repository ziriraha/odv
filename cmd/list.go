package cmd

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
)

// Print the list in the following format:
// ceu - branch -> this branch is present in community, enterprise and upgrade
// c u - branch -> this branch is present in community and upgrade
//  e  - branch -> this branch is present in enterprise only

var listCmd = &cobra.Command{
    Use:   "list",
	Aliases: []string{"ls"},
    Short: "List all branches in the repositories.",
    Long:  "Will list all branches in the specified odoo repositories.",
    Run: func(cmd *cobra.Command, args []string) {
		lib.ForEachRepository(func (i int, repoName string, repository *lib.Repository) error {
			repository.GetBranches() // Pre-fetch branch list
			return nil
		}, true)

		branchPresence := make(map[string]string)
		showVersions, _ := cmd.Flags().GetBool("versions")
		lib.ForEachRepository(func (i int, repoName string, repository *lib.Repository) error {
			letter := repoName[0:1]
			for _, branch := range repository.GetBranches() {
				if !lib.IsVersionBranch(branch) || showVersions {
					presence, ok := branchPresence[branch]
					if !ok { presence = strings.Repeat(" ", len(lib.Repositories)) }
					branchPresence[branch] = presence[:i] + letter + presence[i+1:]
				}
			}
			return nil
		}, false)

		lib.ForEachRepository(func (i int, repoName string, repository *lib.Repository) error {
			letter := repoName[0:1]
			colorizedLetter := repository.Color(letter)
			for branch := range branchPresence {
				branchPresence[branch] = strings.ReplaceAll(branchPresence[branch], letter, colorizedLetter)
			}
			return nil
		}, false)

		sortedBranches := slices.SortedFunc(maps.Keys(branchPresence), func(a, b string) int {
			aVersion := lib.GetVersion(a)
			bVersion := lib.GetVersion(b)
			comparison := strings.Compare(aVersion, bVersion)
			if comparison != 0 { return -comparison
			} else { return strings.Compare(a, b) }
		})

		for _, branch := range sortedBranches { fmt.Printf("%s - %s\n", branchPresence[branch], branch) }
	},
}

func init() {
	listCmd.Flags().BoolP("versions", "v", false, "Show version branches.")
    rootCmd.AddCommand(listCmd)
}
