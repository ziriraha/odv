package cmd

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
	"github.com/ziriraha/odv/views"
)

type updateRepoExtra struct {
	branches      []string
	currentIndex  int
	currentBranch string
}

type branchFetchedMsg struct {
	repoIndex int
}

func fetchNextBranch(repoIndex int, repo *lib.Repository, state *views.RepoOperationState, extra *updateRepoExtra) tea.Cmd {
	if extra.currentIndex >= len(extra.branches) {
		return func() tea.Msg {
			return views.RepoOperationDoneMsg{
				RepoIndex: repoIndex,
				Err:       nil,
				Duration:  time.Since(state.StartTime),
			}
		}
	}

	branch := extra.branches[extra.currentIndex]
	return func() tea.Msg {
		var err error
		if branch == extra.currentBranch {
			err = repo.PullRebase(lib.RemoteOrigin, branch)
		} else {
			err = repo.FetchRefspec(lib.RemoteOrigin, branch)
		}

		if err != nil {
			return views.RepoOperationDoneMsg{
				RepoIndex: repoIndex,
				Err:       fmt.Errorf("failed to fetch %s: %w", branch, err),
				Duration:  time.Since(state.StartTime),
			}
		}

		return branchFetchedMsg{repoIndex: repoIndex}
	}
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all version branches.",
	Long:  "Will fetch (refspec) all version branches in all odoo repositories.",
	Run: func(cmd *cobra.Command, args []string) {
		var states []*views.RepoOperationState
		var extras []*updateRepoExtra
		var repoNames []string
		skipped := make(map[int]bool)

		for _, repoName := range lib.SortedRepoNames {
			repository := lib.Repositories[repoName]

			if repoName == lib.WorkspaceRepo {
				s := views.NewRepoOperationState(repoName)
				idx := len(states)
				states = append(states, &s)
				extras = append(extras, &updateRepoExtra{})
				repoNames = append(repoNames, repoName)
				skipped[idx] = true
				continue
			}

			var versionBranches []string
			for _, branch := range repository.GetBranches() {
				if lib.IsVersionBranch(branch) {
					versionBranches = append(versionBranches, branch)
				}
			}
			if len(versionBranches) == 0 {
				continue
			}
			lib.SortBranches(versionBranches)

			s := views.NewRepoOperationState(repoName)
			states = append(states, &s)
			extras = append(extras, &updateRepoExtra{
				branches:      versionBranches,
				currentIndex:  0,
				currentBranch: repository.GetCurrentBranch(),
			})
			repoNames = append(repoNames, repoName)
		}

		if len(states)-len(skipped) == 0 {
			cmd.Println("No repositories to update.")
			return
		}

		failCount, err := views.RepoBranchSpinnerView{
			Title:          "Updating repositories",
			States:         states,
			SkippedIndices: skipped,
			LaunchOp: func(i int) tea.Cmd {
				return fetchNextBranch(i, lib.Repositories[repoNames[i]], states[i], extras[i])
			},
			OnMsg: func(msg tea.Msg, allStates []*views.RepoOperationState) tea.Cmd {
				if m, ok := msg.(branchFetchedMsg); ok {
					extra := extras[m.repoIndex]
					extra.currentIndex++
					return fetchNextBranch(m.repoIndex, lib.Repositories[repoNames[m.repoIndex]], allStates[m.repoIndex], extra)
				}
				return nil
			},
			RenderRepo: func(i int, state *views.RepoOperationState) string {
				if skipped[i] {
					return fmt.Sprintf("%s %s - skipped (%s)\n",
						views.FaintStyle.Render("⊘"),
						views.RenderRepoName(state.Name),
						views.FaintStyle.Render("no remote found"))
				}
				extra := extras[i]
				switch state.Status {
				case views.StatusInProgress:
					if extra.currentIndex < len(extra.branches) {
						return state.RenderInProgress(fmt.Sprintf("fetching '%s' [%d/%d]", extra.branches[extra.currentIndex], extra.currentIndex+1, len(extra.branches)))
					}
					return state.RenderInProgress("finalizing...")
				case views.StatusDone:
					return state.RenderDone(fmt.Sprintf("%d branches", len(extra.branches)))
				case views.StatusFailed:
					return state.RenderFailed("failed to update")
				}
				return ""
			},
		}.Run()

		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
		if failCount > 0 {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
