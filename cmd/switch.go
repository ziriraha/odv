package cmd

import (
	"fmt"
	"sync"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/internal"
)

var switchCmd = &cobra.Command{
    Use:   "switch <branch>",
    Short: "Switch to an existing branch.",
	Long: "Will switch all three odoo repositories to the specified branch or version.",
	Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
		version := internal.DetectVersion(args[0])
		internal.Debug.Printf("switchCmd: version '%v' was detected", version)
		var repoBranch sync.Map
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			branchName := args[0]
			if !repository.BranchExists(branchName) {
				branchName = version
				if !repository.BranchExists(branchName) { branchName = internal.FallbackBranch }
			}
			repoBranch.Store(repoName, branchName)
		}, true)

		var spinners sync.Map
		ms := internal.NewMultiSpinner()
		defer ms.Close()
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			branchName, _ := repoBranch.Load(repoName)
			text := fmt.Sprintf("[%s] Switching to '%s'", repository.Color(repoName), branchName)
			spinners.Store(repoName, ms.Add(text))
		}, false)
		ms.Start()

		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			branchName, _ := repoBranch.Load(repoName)
			err := repository.SwitchBranch(branchName.(string))
			ls, _ := spinners.Load(repoName)
			spinner := ls.(*internal.LineSpinner)
			if err != nil {
				ms.AddOnClose(func() { internal.Error.Printf("in repository %v: %v", repoName, err) })
				ms.Stop(spinner, color.New(color.FgRed, color.Bold).Sprint("✗"))
			} else {
				ms.Stop(spinner, color.New(color.FgGreen, color.Bold).Sprint("✓"))
			}
		}, true)
	},
}

func init() {
    rootCmd.AddCommand(switchCmd)
}
