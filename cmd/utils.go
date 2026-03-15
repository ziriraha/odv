package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
	"github.com/ziriraha/odv/views"
)

var utilsCmd = &cobra.Command{
	Use:   "utils",
	Short: "Utilities for managing odoo.",
}

func findKillOdooProcess() error {
	odooPort := lib.GetConfig().OdooPort
	pid, err := exec.Command("lsof", "-ti", fmt.Sprintf(":%d", odooPort)).CombinedOutput()
	if err != nil || len(pid) == 0 {
		return fmt.Errorf("no process found listening on port %d", odooPort)
	}
	pidInt, err := strconv.Atoi(strings.TrimSpace(string(pid)))
	if err != nil {
		return fmt.Errorf("invalid pid: %v", err)
	}
	process, err := os.FindProcess(pidInt)
	if err != nil {
		return fmt.Errorf("could not find process: %v", err)
	}
	err = process.Kill()
	if err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}
	return nil
}

var utilsKillOdooCmd = &cobra.Command{
	Use:   "kill-odoo",
	Short: "Find and kill the odoo process.",
	Long:  fmt.Sprintf("Finds the pid of the process listening on port %d and kills it.", lib.GetConfig().OdooPort),
	Run: func(cmd *cobra.Command, args []string) {
		err := findKillOdooProcess()
		if err != nil {
			cmd.PrintErrln("Failed to kill odoo process:", err)
			os.Exit(1)
		}
		cmd.Println("Odoo process killed successfully.")
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
	utilsCmd.AddCommand(utilsKillOdooCmd)
	utilsCmd.AddCommand(utilsCleanBranchesCmd)
	utilsCmd.AddCommand(utilsDeleteBranchCmd)

	rootCmd.AddCommand(utilsCmd)
}
