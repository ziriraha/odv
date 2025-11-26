package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/internal"
)

var updateCmd = &cobra.Command{
    Use:   "update",
    Short: "Update current branch.",
	Long: "Will fetch and pull (ff-only) the current branch in all three odoo repositories.",
    Run: func(cmd *cobra.Command, args []string) {
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			if repoName == ".vscode" { return }
			curBranch, err := repository.GetCurrentBranch()
			if err != nil {
				internal.Error.Printf("in repository %v: getting current branch: %v", repoName, err)
				return
			}
			if !internal.IsVersionBranch(curBranch) {
				internal.Debug.Printf("in repository %v: current branch %v is not a version branch, skipping update", repoName, curBranch)
				return
			}
			err = repository.Fetch(curBranch)
			if err != nil {
				internal.Error.Printf("in repository %v: fetching branch %v: %v", repoName, curBranch, err)
				return
			}
			err = repository.Pull()
			if err != nil { 
				internal.Error.Printf("in repository %v: pulling branch %v: %v", repoName, curBranch, err) 
			}
		}, true)
	},
}

func init() {
    rootCmd.AddCommand(updateCmd)
}
