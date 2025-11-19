package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/internal"
)

var statusCmd = &cobra.Command{
    Use:   "status",
    Short: "Prints current branch.",
	Long: "Will print the current branch in all three odoo repositories.",
    Run: func(cmd *cobra.Command, args []string) {
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			curBranch, err := repository.GetCurrentBranch()
			if err != nil { internal.Error.Printf("in repository %v: getting current branch: %v", repoName, err)
			} else { fmt.Printf("%v: %v\n", repository.Color(repoName), curBranch) }
		}, false)
	},
}

func init() {
    rootCmd.AddCommand(statusCmd)
}
