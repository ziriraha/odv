package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odoodev/internal"
)

var updateCmd = &cobra.Command{
    Use:   "update",
    Short: "Update current branch.",
	Long: "Will fetch and pull (ff-only) the current branch in all three odoo repositories.",
    Run: func(cmd *cobra.Command, args []string) {
		internal.ForEachRepository(func (repository *internal.Repository) error {
			curBranch, err := repository.GetCurrentBranch()
			if err != nil {
				return fmt.Errorf("getting current branch: %w", err)
			}
			err = repository.Fetch(curBranch)
			if err != nil {
				return fmt.Errorf("fetching branch %v: %w", curBranch, err)
			}
			err = repository.Pull()
			if err != nil {
				return fmt.Errorf("pulling branch %v: %w", curBranch, err)
			}
			return nil
		})
	},
}

func init() {
    rootCmd.AddCommand(updateCmd)
}
