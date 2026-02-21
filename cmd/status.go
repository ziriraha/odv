package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
	"github.com/ziriraha/odv/views"
)

func colorizeIndicator(s string) string {
	switch s {
	case "A", "R":
		return views.DiffAddedStyle.Render(s)
	case "M":
		return views.DiffModifiedStyle.Render(s)
	case "D":
		return views.DiffDeletedStyle.Render(s)
	case "?":
		return views.DiffUntrackedStyle.Render(s)
	default:
		return s
	}
}

func colorizeStatusIndicator(status string) string {
	parts := strings.Split(status, "")
	if len(parts) != 2 {
		return status
	}

	switch status {
	case "UU", "AA", "DD", "AU", "UA", "DU", "UD":
		return views.DiffConflictStyle.Render(status)
	case "!!":
		return views.DiffUntrackedStyle.Render(status)
	default:
		return colorizeIndicator(parts[0]) + colorizeIndicator(parts[1])
	}
}

var statusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{"st"},
	Short:   "Prints current branch's status.",
	Long:    "Will print the current branch in all three odoo repositories.",
	Run: func(cmd *cobra.Command, args []string) {
		short, _ := cmd.Flags().GetBool("short")

		type repoStatus struct {
			name   string
			status string
		}

		statuses := make([]repoStatus, len(lib.SortedRepoNames))
		var wg sync.WaitGroup

		for i, repoName := range lib.SortedRepoNames {
			wg.Go(func() {
				repository := lib.Repositories[repoName]
				var repoWg sync.WaitGroup
				var ahead, behind int
				var changes []string

				curBranch := repository.GetCurrentBranch()
				repoWg.Go(func() {
					remote := lib.RemoteOrigin
					if repoName != "upgrade" {
						remote = lib.GetRemoteForBranch(curBranch)
					}
					ahead, behind, _ = repository.GetAheadBehindInfo(remote, curBranch)
				})
				repoWg.Go(func() { changes, _ = repository.GetStatus() })
				repoWg.Wait()

				output := strings.Builder{}
				if ahead < 0 && behind < 0 {
					curBranch = views.LocalBranchStyle.Render(curBranch)
				}
				output.WriteString(fmt.Sprintf("%s %s - %s ", views.FaintStyle.Render("*"), views.RenderRepoName(repoName), curBranch))
				if ahead > 0 {
					output.WriteString(views.AheadStyle.Render(fmt.Sprintf("↑%d", ahead)))
				}
				if behind > 0 {
					output.WriteString(views.BehindStyle.Render(fmt.Sprintf("↓%d", behind)))
				}
				output.WriteString("\n")
				if !short {
					for _, change := range changes {
						indicator := colorizeStatusIndicator(change[0:2])
						change = fmt.Sprintf("%s %s", indicator, change[3:])
						output.WriteString(fmt.Sprintf("   |%s\n", change))
					}
				}
				statuses[i] = repoStatus{name: repoName, status: output.String()}
			})
		}
		wg.Wait()

		for _, status := range statuses {
			fmt.Print(status.status)
		}
	},
}

func init() {
	statusCmd.Flags().BoolP("short", "s", false, "Do not show changes (shorter version).")
	rootCmd.AddCommand(statusCmd)
}
