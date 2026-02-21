package cmd

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
	"github.com/ziriraha/odv/views"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List available branches.",
	Long:    "Will list all branches in the specified odoo repositories with color-coded presence indicators.",
	Run: func(cmd *cobra.Command, args []string) {
		showVersions, _ := cmd.Flags().GetBool("all")
		branchPresence := make(map[string][]bool)

		for _, branch := range lib.GetAllBranches() {
			if !lib.IsVersionBranch(branch) || showVersions {
				branchPresence[branch] = make([]bool, len(lib.SortedRepoNames))
			}
		}

		var wg sync.WaitGroup
		for repoIndex, repoName := range lib.SortedRepoNames {
			repository := lib.Repositories[repoName]
			wg.Go(func() {
				for _, branch := range repository.GetBranches() {
					if _, ok := branchPresence[branch]; ok {
						branchPresence[branch][repoIndex] = true
					}
				}
			})
		}
		wg.Wait()

		branches := slices.Collect(maps.Keys(branchPresence))
		lib.SortBranches(branches)
		for _, branch := range branches {
			var indicator strings.Builder

			for repoIndex, repoName := range lib.SortedRepoNames {
				if branchPresence[branch][repoIndex] {
					indicator.WriteString(views.RenderRepoLetter(repoName))
				} else {
					indicator.WriteString(views.FaintStyle.Render("·"))
				}
			}
			fmt.Printf("%s - %s\n", indicator.String(), branch)
		}
	},
}

func init() {
	listCmd.Flags().BoolP("all", "a", false, "Show all branches.")
	rootCmd.AddCommand(listCmd)
}
