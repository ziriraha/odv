package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ziriraha/odoodev/internal"
)

var updateCmd = &cobra.Command{
    Use:   "update",
    Short: "Update current branch.",
	Long: "Will fetch and pull (ff-only) the current branch in all three odoo repositories.",
    Run: func(cmd *cobra.Command, args []string) {
		internal.ForEachRepository(func (i int, repository *internal.Repository) {
			curBranch, err := repository.GetCurrentBranch()
			if err != nil { internal.Error.Printf("in repository %v: getting current branch: %v", repository.Name, err)
			} else {
				err = repository.Fetch(curBranch)
				if err != nil { internal.Error.Printf("in repository %v: fetching branch %v: %v", repository.Name, curBranch, err)
				} else {
					err = repository.Pull()
					if err != nil { internal.Error.Printf("in repository %v: pulling branch %v: %v", repository.Name, curBranch, err) }
				}
			}
		}, true)
	},
}

func init() {
    rootCmd.AddCommand(updateCmd)
}
