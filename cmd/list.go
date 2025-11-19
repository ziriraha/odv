package cmd

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/internal"
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
		var branches sync.Map
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			branches.Store(repoName, repository.GetBranches())
		}, true)

		branchPresence := make(map[string]string)
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			branchList, _ := branches.Load(repoName)
			letter := repoName[0:1]
			for _, branch := range branchList.([]string) {
				presence, ok := branchPresence[branch]
				if !ok { presence = strings.Repeat(" ", len(internal.Repositories)) }
				presence = presence[:i] + letter + presence[i+1:]
				branchPresence[branch] = presence
			}
		}, false)

		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			letter := repoName[0:1]
			colorizedLetter := repository.Color(letter)
			for branch := range branchPresence {
				branchPresence[branch] = strings.ReplaceAll(branchPresence[branch], letter, colorizedLetter)
			}
		}, false)

		for _, branch := range slices.Sorted(maps.Keys(branchPresence)) {
			fmt.Printf("%s - %s\n", branchPresence[branch], branch)
		}
	},
}

func init() {
    rootCmd.AddCommand(listCmd)
}
