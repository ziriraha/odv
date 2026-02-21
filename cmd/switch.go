package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
	"github.com/ziriraha/odv/views"
)

func performSwitch(repoIndex int, repo *lib.Repository, state *views.RepoOperationState, targetBranch string) tea.Cmd {
	return func() tea.Msg {
		startTime := time.Now()

		if state.Name == lib.WorkspaceRepo {
			curBranch := repo.GetCurrentBranch()
			if curBranch != "" && curBranch != lib.FallbackBranch {
				changes, err := repo.GetStatus()
				if err == nil && len(changes) > 0 {
					commitMessage := fmt.Sprintf("odv auto-commit %v\n\nBefore switching to '%s'", time.Now().Format(time.RFC3339), targetBranch)
					if err := repo.CommitAll(commitMessage); err != nil {
						return views.RepoOperationDoneMsg{
							RepoIndex: repoIndex,
							Err:       fmt.Errorf("auto-commit failed before switch: %w", err),
							Duration:  time.Since(startTime),
						}
					}
				}
			}
		}

		return views.RepoOperationDoneMsg{
			RepoIndex: repoIndex,
			Err:       repo.SwitchBranch(targetBranch),
			Duration:  time.Since(startTime),
		}
	}
}

var switchCmd = &cobra.Command{
	Use:   "switch [branch]",
	Short: "Switch to an existing branch.",
	Long:  "If a branch is specified, switch to it directly. If no branch is specified, displays a list to choose from.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var selectedBranch string
		branches := lib.GetAllBranches()
		if len(branches) == 0 {
			fmt.Println("No branches found.")
			os.Exit(1)
		}

		if len(args) == 0 {
			choice, err := views.BranchSelectListView{
				Title:    "Select a branch to switch to",
				Branches: branches,
			}.Run()

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
				os.Exit(1)
			}
			if choice == "" {
				return
			}
			selectedBranch = choice
		} else {
			selectedBranch, _ = strings.CutPrefix(args[0], "odoo-dev:")
		}

		if !slices.Contains(branches, selectedBranch) {
			fmt.Fprintf(os.Stderr, "branch '%s' was not found\n", selectedBranch)
			os.Exit(1)
		}
		version := lib.DetectVersion(selectedBranch)

		repoBranches := make(map[string]string)
		for _, repoName := range lib.SortedRepoNames {
			repository := lib.Repositories[repoName]
			branchName := selectedBranch
			if !repository.BranchExists(branchName) {
				branchName = version
				if !repository.BranchExists(branchName) {
					branchName = lib.FallbackBranch
					if !repository.BranchExists(branchName) {
						fmt.Fprintf(os.Stderr, "no suitable branch found for '%s' in repo '%s' (tried: %s, %s, %s)\n", selectedBranch, repoName, selectedBranch, version, lib.FallbackBranch)
						os.Exit(1)
					}
				}
			}
			repoBranches[repoName] = branchName
		}

		states := make([]*views.RepoOperationState, len(lib.SortedRepoNames))
		targetBranches := make([]string, len(lib.SortedRepoNames))
		for i, repoName := range lib.SortedRepoNames {
			s := views.NewRepoOperationState(repoName)
			states[i] = &s
			targetBranches[i] = repoBranches[repoName]
		}

		views.RepoBranchSpinnerView{
			Title:  "Switching branches",
			States: states,
			LaunchOp: func(i int) tea.Cmd {
				return performSwitch(i, lib.Repositories[lib.SortedRepoNames[i]], states[i], targetBranches[i])
			},
			RenderRepo: func(i int, state *views.RepoOperationState) string {
				tb := targetBranches[i]
				switch state.Status {
				case views.StatusInProgress:
					return state.RenderInProgress(fmt.Sprintf("switching to '%s'", tb))
				case views.StatusDone:
					return state.RenderDone(fmt.Sprintf("switched to '%s'", tb))
				case views.StatusFailed:
					return state.RenderFailed(fmt.Sprintf("failed to switch to '%s'", tb))
				}
				return ""
			},
		}.RunOrExit()
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
