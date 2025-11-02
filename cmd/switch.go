package cmd

import (
	"fmt"
	"log"
	"sync"

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
		var wg sync.WaitGroup
		for _, repository := range internal.GetRepositories() {
			wg.Add(1)
			go func(repository *internal.Repository) {
				defer wg.Done()
				branchName := findBranch(repository, args[0])
				fmt.Printf("Switching %v to branch %v\n", repository.Name, branchName)
				err := repository.SwitchBranch(branchName)
				if err != nil {
					log.Fatal(fmt.Errorf("switching to branch %v in repository %v: %w", branchName, repository.Name, err))
				}
			}(&repository)
		}
		wg.Wait()
    },
}

func init() {
    rootCmd.AddCommand(switchCmd)
}
