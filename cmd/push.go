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

type pushRepoExtra struct {
	remote     string
	branch     string
	skipReason string
}

func performPush(repoIndex int, repo *lib.Repository, extra *pushRepoExtra, force bool) tea.Cmd {
	return func() tea.Msg {
		startTime := time.Now()
		return views.RepoOperationDoneMsg{
			RepoIndex: repoIndex,
			Err:       repo.Push(extra.remote, extra.branch, force),
			Duration:  time.Since(startTime),
		}
	}
}

var pushForce bool

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes current branch to dev remote.",
	Long:  "Pushes the current branch in all repositories to the dev remote. Skips version branches.",
	Run: func(cmd *cobra.Command, args []string) {
		var states []*views.RepoOperationState
		var extras []*pushRepoExtra
		var repoNames []string
		skipped := make(map[int]bool)

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
			extra := &pushRepoExtra{remote: remote, branch: curBranch}
			idx := len(states)

			if lib.IsVersionBranch(curBranch) {
				extra.skipReason = "on version branch"
				skipped[idx] = true
			}

			states = append(states, &s)
			extras = append(extras, extra)
			repoNames = append(repoNames, repoName)
		}

		if len(states)-len(skipped) == 0 {
			cmd.Println("No repositories on non-version branches to push.")
			return
		}

		failCount, err := views.RepoBranchSpinnerView{
			Title:          "Pushing branches",
			States:         states,
			SkippedIndices: skipped,
			LaunchOp: func(i int) tea.Cmd {
				return performPush(i, lib.GetRepository(states[i].Name), extras[i], pushForce)
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
					return state.RenderInProgress(fmt.Sprintf("pushing '%s'", extra.branch))
				case views.StatusDone:
					return state.RenderDone(fmt.Sprintf("pushed '%s'", extra.branch))
				case views.StatusFailed:
					return state.RenderFailed(fmt.Sprintf("failed to push '%s'", extra.branch))
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
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().BoolVar(&pushForce, "force", false, "Force push with --force-with-lease.")
}
