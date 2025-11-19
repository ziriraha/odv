package cmd

import (
	"fmt"

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
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			branchName := args[0]
			if !repository.BranchExists(branchName) {
				branchName = version
				if !repository.BranchExists(branchName) { branchName = repository.DefaultBranch }
			}
			fmt.Printf("Switching '%v' to branch '%v'\n", repository.Color(repoName), branchName)
			err := repository.SwitchBranch(branchName)
			if err != nil { 
				internal.Error.Printf("in repository %v: switching to branch '%v': %v", repoName, branchName, err) 
			}
		}, true)
	},
}

func init() {
    rootCmd.AddCommand(switchCmd)
}
