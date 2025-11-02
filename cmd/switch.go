package cmd

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/ziriraha/odoodev/internal"
)

func findBranch(repository *git.Repository, branchName string) string {
	_, err := repository.Branch(branchName)
	if err != nil {
		version := internal.DetectVersion(branchName)
		_, err := repository.Branch(version)
		if err != nil {
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
		for name, repository := range internal.GetRepositories() {
			wg.Add(1)
			go func(name string, repository *git.Repository) {
				defer wg.Done()
				branchName := findBranch(repository, args[0])
				fmt.Printf("Switching %v to branch %v\n", name, branchName)
				err := internal.SwitchBranch(repository, branchName)
				if err != nil {
					log.Fatal(fmt.Errorf("switching to branch %v in repository %v: %w", branchName, name, err))
				}
			}(name, repository)
		}
		wg.Wait()
    },
}

func init() {
    rootCmd.AddCommand(switchCmd)
}
