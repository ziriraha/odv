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

type pullRepoExtra struct {
	branch     string
	skipReason string
}

func performPull(repoIndex int, repo *lib.Repository, extra *pullRepoExtra) tea.Cmd {
	return func() tea.Msg {
		startTime := time.Now()
		return views.RepoOperationDoneMsg{
			RepoIndex: repoIndex,
			Err:       repo.PullRebase(lib.GetRemoteForBranch(extra.branch), extra.branch),
			Duration:  time.Since(startTime),
		}
	}
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls current branch.",
	Long:  "Will pull (ff-only) the current branch in all three odoo repositories.",
	Run: func(cmd *cobra.Command, args []string) {
		var states []*views.RepoOperationState
		var extras []*pullRepoExtra
		var repoNames []string
		skipped := make(map[int]bool)

		for _, repoName := range lib.SortedRepoNames {
			repository := lib.Repositories[repoName]
			curBranch := repository.GetCurrentBranch()
			s := views.NewRepoOperationState(repoName)

			extra := &pullRepoExtra{branch: curBranch}
			idx := len(states)

			if repoName == lib.WorkspaceRepo {
				extra.skipReason = "no remote found"
				skipped[idx] = true
			} else if !lib.IsVersionBranch(curBranch) {
				extra.skipReason = "not on version branch"
				skipped[idx] = true
			}

			states = append(states, &s)
			extras = append(extras, extra)
			repoNames = append(repoNames, repoName)
		}

		if len(states)-len(skipped) == 0 {
			cmd.Println("No repositories on version branches to pull.")
			return
		}

		failCount, err := views.RepoBranchSpinnerView{
			Title:          "Pulling branches",
			States:         states,
			SkippedIndices: skipped,
			LaunchOp: func(i int) tea.Cmd {
				return performPull(i, lib.Repositories[repoNames[i]], extras[i])
			},
			RenderRepo: func(i int, state *views.RepoOperationState) string {
				extra := extras[i]
				if skipped[i] {
					return fmt.Sprintf("%s %s - skipped (%s)\n",
						views.FaintStyle.Render("⊘"),
						views.RenderRepoName(state.Name),
						views.FaintStyle.Render(extra.skipReason))
				}
				switch state.Status {
				case views.StatusInProgress:
					return state.RenderInProgress(fmt.Sprintf("pulling '%s'", extra.branch))
				case views.StatusDone:
					return state.RenderDone(fmt.Sprintf("pulled '%s'", extra.branch))
				case views.StatusFailed:
					return state.RenderFailed(fmt.Sprintf("failed to pull '%s'", extra.branch))
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
	rootCmd.AddCommand(pullCmd)
}
