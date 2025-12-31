package cmd

import (
	"fmt"
	"sync"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/internal"
)

var switchCmd = &cobra.Command{
    Use:   "switch <branch>",
    Short: "Switch to an existing branch.",
	Long: "Will switch all three odoo repositories to the specified branch or version.",
	Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
    	input, _ := strings.CutPrefix(args[0], "odoo-dev:")
		version := internal.DetectVersion(input)
		internal.Debug.Printf("switchCmd: version '%v' was detected", version)
		var repoBranch sync.Map
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) error {
			branchName := input
			if !repository.BranchExists(branchName) {
				branchName = version
				if !repository.BranchExists(branchName) { branchName = internal.FallbackBranch }
			}
			repoBranch.Store(repoName, branchName)
			return nil
		}, true)

		var spinners sync.Map
		ms := internal.NewMultiSpinner()
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) error {
			branchName, _ := repoBranch.Load(repoName)
			text := fmt.Sprintf("[%s] Switching to '%s'", repository.Color(repoName), branchName)
			spinners.Store(repoName, ms.Add(text))
			return nil
		}, false)
		ms.Start()

		errors := internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) error {
			branchName, _ := repoBranch.Load(repoName)
			err := repository.SwitchBranch(branchName.(string))
			ls, _ := spinners.Load(repoName)
			spinner := ls.(*internal.LineSpinner)
			if err != nil {
				ms.Fail(spinner)
			} else {
				ms.Done(spinner)
			}
			return err
		}, true)
		ms.Close()
		internal.PrintRepositoryErrors(errors)
	},
}

func init() {
    rootCmd.AddCommand(switchCmd)
}
