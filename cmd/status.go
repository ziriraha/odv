package cmd

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
)

// The format will be:
// [repoName] currentBranch ↑3↓2
//     ...changes

func colorizeIndicator(s string) string {
	switch s {
	case "A", "R": return color.GreenString(s)
	case "M": return color.YellowString(s)
	case "D": return color.RedString(s)
	case "?": return lib.GrayString(s)
	default: return string(s)
	}
}

func colorizeStatusIndicator(status string) string {
	parts := strings.Split(status, "")
	if len(parts) != 2 { return status }

	switch status {
	case "UU", "AA", "DD", "AU", "UA", "DU", "UD": return color.New(color.FgRed, color.Bold).Sprint(status)
	case "!!": return lib.GrayString(status)
	default: return colorizeIndicator(parts[0]) + colorizeIndicator(parts[1])
	}
}

var statusCmd = &cobra.Command{
    Use:   "status",
	Aliases: []string{"st"},
    Short: "Prints current branch.",
	Long: "Will print the current branch in all three odoo repositories.",
    Run: func(cmd *cobra.Command, args []string) {
		var statuses sync.Map
		short, _ := cmd.Flags().GetBool("short")
		lib.ForEachRepository(func (i int, repoName string, repository *lib.Repository) error {
			curBranch := repository.GetCurrentBranch()
			remote := "dev"
			if lib.IsVersionBranch(curBranch) || repoName == "upgrade" { remote = "origin" }
			ahead, behind, err := repository.GetAheadBehindInfo(remote, curBranch)
			if err != nil {
				lib.Debug.Printf("in repository %v: error getting ahead/behind info probably due to no upstream branch: %v", repoName, err)
			}
			changes, err := repository.GetStatus()
			if err != nil {
				lib.Debug.Printf("in repository %v: error getting changes, won't show them: %v", repoName, err)
				changes = []string{}
			}

			output := strings.Builder{}
			if ahead < 0 && behind < 0 { curBranch = lib.GrayString(curBranch) }
			output.WriteString(fmt.Sprintf("[%s] %s ", repository.Color(repoName), curBranch))
			if ahead > 0 { output.WriteString(color.GreenString("↑%d", ahead)) }
			if behind > 0 { output.WriteString(color.RedString("↓%d", behind)) }
			output.WriteString("\n")
			if !short {
				for _, change := range changes {
					indicator := colorizeStatusIndicator(change[0:2])
					change = fmt.Sprintf("%s %s", indicator, change[3:])
					output.WriteString(fmt.Sprintf("   |%s\n", change))
				}
			}
			statuses.Store(repoName, output.String())
			return nil
		}, true)

		lib.ForEachRepository(func (i int, repoName string, repository *lib.Repository) error {
			if statusStr, ok := statuses.Load(repoName); ok { fmt.Print(statusStr.(string)) }
			return nil
		}, false)
	},
}

func init() {
	statusCmd.Flags().BoolP("short", "s", false, "Do not show changes (shorter version).")
    rootCmd.AddCommand(statusCmd)
}
