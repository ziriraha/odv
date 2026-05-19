package cmd

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
	"github.com/ziriraha/odv/views"
)

var utilsCmd = &cobra.Command{
	Use:   "utils",
	Short: "Utilities for managing odoo.",
}

var utilsKillCmd = &cobra.Command{
	Use:   "kill [port]",
	Short: "Find and kill process by port (default 8069).",
	Long:  "Finds the pid of the process listening on port (default 8069) and kills it.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		port := 8069
		if len(args) == 1 {
			port, err = strconv.Atoi(args[0])
			if err != nil {
				cmd.PrintErrln("Invalid port number:", err)
				os.Exit(1)
			}
		}
		err = lib.KillProcessByPort(port)
		if err != nil {
			cmd.PrintErrln("Failed to kill process:", err)
			os.Exit(1)
		}
		cmd.Printf("Process running in %d killed successfully.\n", port)
	},
}

var utilsCleanBranchesCmd = &cobra.Command{
	Use:   "clean-branches",
	Short: "Clean up local .workspace git branches that have been deleted in other repos.",
	Run: func(cmd *cobra.Command, args []string) {
		workspaceRepo := lib.GetRepository(".workspace")
		branchesToKeep := make(map[string]struct{})
		branchesToKeep["main"] = struct{}{} // always keep main branch
		for _, branch := range lib.GetAllBranches() {
			branchesToKeep[branch] = struct{}{}
		}

		var deletedCount int
		for _, branch := range workspaceRepo.GetBranches() {
			if _, exists := branchesToKeep[branch]; !exists {
				err := workspaceRepo.DeleteBranch(branch)
				if err != nil {
					cmd.PrintErrf("Failed to delete branch '%s': %v\n", branch, err)
				} else {
					cmd.Printf("Deleted orphaned branch '%s'\n", branch)
					deletedCount++
				}
			}
		}
	},
}

var utilsDeleteBranchCmd = &cobra.Command{
	Use:   "delete-branch <branch>",
	Short: "Delete the specified branch in all repositories.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		branchToDelete := args[0]
		for repoName, repository := range lib.GetRepositories() {
			if repository.BranchExists(branchToDelete) {
				err := repository.DeleteBranch(branchToDelete)
				if err != nil {
					cmd.PrintErrln(views.RepoLine(repoName, "Failed to delete branch '%s': %v", branchToDelete, err))
				} else {
					cmd.Println(views.RepoLine(repoName, "Deleted branch '%s'", branchToDelete))
				}
			}
		}
	},
}

func init() {
	utilsCmd.AddCommand(utilsKillCmd)
	utilsCmd.AddCommand(utilsCleanBranchesCmd)
	utilsCmd.AddCommand(utilsDeleteBranchCmd)

	rootCmd.AddCommand(utilsCmd)
}
