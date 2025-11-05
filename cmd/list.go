package cmd

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odoodev/internal"
)

// Print the list in the following format:
// ceu - branch -> this branch is present in community, enterprise and upgrade
// c u - branch -> this branch is present in community and upgrade
//  e   - branch -> this branch is present in enterprise only
// Sort the output branches by most present to least present and then alphabetically

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List all branches in the repositories.",
    Long:  "Will list all branches in the specified odoo repositories.",
    Run: func(cmd *cobra.Command, args []string) {
		branches := make(map[*internal.Repository][]string)
		letters := make([]string, 0, len(internal.Repositories))
		var mapLock sync.Mutex
		internal.ForEachRepository(func (repository *internal.Repository) error {
			branchList, err := repository.GetBranches()
			letters = append(letters, repository.Name[0:1])
			if err != nil {
				return fmt.Errorf("getting branches: %w", err)
			}
			mapLock.Lock()
			branches[repository] = branchList
			mapLock.Unlock()
			return nil
		}, true)
		sort.Strings(letters)

		branchPresence := make(map[string]string)
		for repo, branchList := range branches {
			for _, branch := range branchList {
				presence, ok := branchPresence[branch]
				if !ok {
					presence = strings.Repeat(" ", len(branches))
				}
				letter := repo.Name[0:1]

				if pos := sort.SearchStrings(letters, letter); pos == 0 {
					presence = letter + presence[1:]
				} else if pos == len(letters)-1 {
					presence = presence[:pos] + letter
				} else {
					presence = presence[:pos] + letter + presence[pos+1:]
				}
				branchPresence[branch] = presence
			}
		}

		type branchInfo struct {
			name     string
			presence string
		}
		sortedBranches := make([]branchInfo, 0, len(branchPresence))
		for branch, presence := range branchPresence {
			sortedBranches = append(sortedBranches, branchInfo{
				name:     branch,
				presence: presence,
			})
		}
		sort.Slice(sortedBranches, func(i, j int) bool {
			iLen := len(strings.TrimSpace(sortedBranches[i].presence))
			jLen := len(strings.TrimSpace(sortedBranches[j].presence))
			if iLen != jLen {
				return iLen > jLen
			}
			return sortedBranches[i].name < sortedBranches[j].name
		})
		for _, branch := range sortedBranches {
			internal.ForEachRepository(func (r *internal.Repository) error {
				letter := r.Name[0:1]
				branch.presence = strings.ReplaceAll(branch.presence, letter, r.Color(letter))
				return nil
			}, false)
			fmt.Printf("%s - %s\n", branch.presence, branch.name)
		}
	},
}

func init() {
    rootCmd.AddCommand(listCmd)
}
