package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odoodev/internal"
)

func findBranch(repository *internal.Repository, branchName string) string {
	if !repository.BranchExists(branchName) {
		version := internal.DetectVersion(branchName)
		if !repository.BranchExists(version) {
			return "master"
		}
		return version
	}
	return branchName
}

var switchCmd = &cobra.Command{
    Use:   "switch [branch]",
    Short: "Switch to an existing branch.",
	Long: "Will switch all three odoo repositories to the specified branch or version.",
	Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
		internal.ForEachRepository(func (repository *internal.Repository) error {
			branchName := findBranch(repository, args[0])
			fmt.Printf("Switching %v to branch %v\n", repository.Color(repository.Name), branchName)
			err := repository.SwitchBranch(branchName)
			if err != nil {
				return fmt.Errorf("switching to branch %v: %w", branchName, err)
			}
			return nil
		})
	},
}

func init() {
    rootCmd.AddCommand(switchCmd)
}
