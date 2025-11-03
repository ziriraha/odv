package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/ziriraha/odoodev/internal"
)

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List all branches in the repositories.",
    Long:  "Will list all branches in the specified odoo repositories.",
    Run: func(cmd *cobra.Command, args []string) {
		branches := make(map[string][]string)
		internal.ForEachRepository(func (repository *internal.Repository) error {
			branchList, err := repository.GetBranches()
			if err != nil {
				return fmt.Errorf("getting branches: %w", err)
			}
			branches[repository.Name] = branchList
			return nil
		})

		// Print the list in the following format:
		// ceu - branch -> this branch is present in community, enterprise and upgrade
		// c u - branch -> this branch is present in community and upgrade
		//  e   - branch -> this branch is present in enterprise only
		branchPresence := make(map[string]string)
		for repoName, branchList := range branches {
			for _, branch := range branchList {
				presence, ok := branchPresence[branch]
				if !ok {
					presence = "   "
				}
				switch repoName {
				case "community":
					presence = "c" + presence[1:]
				case "enterprise":
					presence = presence[:1] + "e" + presence[2:]
				case "upgrade":
					presence = presence[:2] + "u"
				}
				branchPresence[branch] = presence
			}
		}

		// Sort the branches by most present to least present and then alphabetically
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
			// Colorize presence
			branch.presence = strings.ReplaceAll(branch.presence, "c", color.GreenString("c"))
			branch.presence = strings.ReplaceAll(branch.presence, "e", color.YellowString("e"))
			branch.presence = strings.ReplaceAll(branch.presence, "u", color.CyanString("u"))
			fmt.Printf("%s - %s\n", branch.presence, branch.name)
		}
	},
}

func init() {
    rootCmd.AddCommand(listCmd)
}
