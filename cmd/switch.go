package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/internal"
)

var switchCmd = &cobra.Command{
    Use:   "switch <branch>",
    Short: "Switch to an existing branch.",
	Long: "Will switch all three odoo repositories to the specified branch or version.",
	Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
		version := internal.DetectVersion(args[0])
		internal.Debug.Printf("switchCmd: version '%v' was detected", version)
		var repoBranch sync.Map
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			branchName := args[0]
			if !repository.BranchExists(branchName) {
				branchName = version
				if !repository.BranchExists(branchName) { branchName = repository.DefaultBranch }
			}
			repoBranch.Store(repoName, branchName)
		}, true)

		var wg sync.WaitGroup
		wg.Go(func() {
			internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
				branchName, _ := repoBranch.Load(repoName)
				fmt.Printf("[%s] Switching to '%s'\n", repository.Color(repoName), branchName)
			}, false)
		})

		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			branchName, _ := repoBranch.Load(repoName)
			err := repository.SwitchBranch(branchName.(string))
			wg.Wait()
			if err != nil {
				internal.Error.Printf("in repository %v: %v", repoName, err) 
			}
		}, true)
	},
}

func init() {
    rootCmd.AddCommand(switchCmd)
}
