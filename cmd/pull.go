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
	remote string
	branch string
}

func performPull(repoIndex int, repo *lib.Repository, extra *pullRepoExtra) tea.Cmd {
	return func() tea.Msg {
		startTime := time.Now()
		return views.RepoOperationDoneMsg{
			RepoIndex: repoIndex,
			Err:       repo.PullRebase(extra.remote, extra.branch),
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

		for _, repoName := range lib.GetSortedRepoNames() {
			if repoName == ".workspace" {
				continue
			}
			repository := lib.GetRepository(repoName)
			curBranch := repository.GetCurrentBranch()
			s := views.NewRepoOperationState(repoName)

			remote := lib.GetRemoteForBranch(curBranch)
			if repoName == "upgrade" {
				remote = lib.RemoteOrigin
			}

			extra := &pullRepoExtra{remote: remote, branch: curBranch}

			states = append(states, &s)
			extras = append(extras, extra)
			repoNames = append(repoNames, repoName)
		}

		failCount, err := views.RepoBranchSpinnerView{
			Title:  "Pulling branches",
			States: states,
			LaunchOp: func(i int) tea.Cmd {
				return performPull(i, lib.GetRepository(states[i].Name), extras[i])
			},
			RenderRepo: func(i int, state *views.RepoOperationState) string {
				extra := extras[i]
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
