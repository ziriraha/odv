package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
	"github.com/ziriraha/odv/views"
)

type rebaseRepoExtra struct {
	branch     string
	skipReason string
	conflicts  []string
}

func isConflictIndicator(status string) bool {
	switch status {
	case "UU", "AA", "DD", "AU", "UA", "DU", "UD":
		return true
	}
	return false
}

func performRebase(repoIndex int, repo *lib.Repository, extra *rebaseRepoExtra) tea.Cmd {
	return func() tea.Msg {
		startTime := time.Now()
		err := repo.PullRebase(lib.RemoteOrigin, extra.branch)

		if err != nil {
			changes, _ := repo.GetStatus()
			for _, change := range changes {
				if len(change) >= 2 && isConflictIndicator(change[0:2]) {
					extra.conflicts = append(extra.conflicts, change)
				}
			}
		}

		return views.RepoOperationDoneMsg{
			RepoIndex: repoIndex,
			Err:       err,
			Duration:  time.Since(startTime),
		}
	}
}

var rebaseCmd = &cobra.Command{
	Use:   "rebase",
	Short: "Rebase current branch on its version branch.",
	Long:  "Will run git pull --rebase origin <versionBranch> on all repositories. Conflicts are left for the user to resolve.",
	Run: func(cmd *cobra.Command, args []string) {
		var states []*views.RepoOperationState
		var extras []*rebaseRepoExtra
		var repoNames []string
		skipped := make(map[int]bool)

		for _, repoName := range lib.SortedRepoNames {
			repository := lib.Repositories[repoName]
			curBranch := repository.GetCurrentBranch()
			s := views.NewRepoOperationState(repoName)

			version := lib.DetectVersion(curBranch)
			extra := &rebaseRepoExtra{branch: version}
			idx := len(states)

			if repoName == lib.WorkspaceRepo {
				extra.skipReason = "no remote found"
				skipped[idx] = true
			} else if curBranch == version {
				extra.skipReason = "already on that base"
				skipped[idx] = true
			}

			states = append(states, &s)
			extras = append(extras, extra)
			repoNames = append(repoNames, repoName)
		}

		if len(states)-len(skipped) == 0 {
			cmd.Println("Nothing to rebase.")
			return
		}

		_, err := views.RepoBranchSpinnerView{
			Title:          "Rebasing branches",
			States:         states,
			SkippedIndices: skipped,
			LaunchOp: func(i int) tea.Cmd {
				return performRebase(i, lib.Repositories[repoNames[i]], extras[i])
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
					return state.RenderInProgress(fmt.Sprintf("rebasing on '%s'", extra.branch))
				case views.StatusDone:
					return state.RenderDone(fmt.Sprintf("rebased on '%s'", extra.branch))
				case views.StatusFailed:
					hasConflicts := len(extra.conflicts) > 0
					if hasConflicts {
						var b strings.Builder
						b.WriteString(fmt.Sprintf("%s %s - conflicts rebasing on '%s'\n",
							views.Cross,
							views.RenderRepoName(state.Name),
							extra.branch))
						for _, change := range extra.conflicts {
							indicator := views.ColorizeStatusIndicator(change[0:2])
							b.WriteString(fmt.Sprintf("   |%s %s\n", indicator, change[3:]))
						}
						return b.String()
					}
					return state.RenderFailed(fmt.Sprintf("failed to rebase on '%s'", extra.branch))
				}
				return ""
			},
		}.Run()

		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(rebaseCmd)
}
