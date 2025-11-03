package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odoodev/internal"
)

var statusCmd = &cobra.Command{
    Use:   "status",
    Short: "Prints current branch.",
	Long: "Will print the current branch in all three odoo repositories.",
    Run: func(cmd *cobra.Command, args []string) {
		internal.ForEachRepository(func (repository *internal.Repository) error {
			curBranch, err := repository.GetCurrentBranch()
			if err != nil {
				return fmt.Errorf("getting current branch: %w", err)
			}
			fmt.Printf("%v: %v\n", repository.Name, curBranch)
			return nil
		})
	},
}

func init() {
    rootCmd.AddCommand(statusCmd)
}
