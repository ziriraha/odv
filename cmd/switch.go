package cmd

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
)

var switchCmd = &cobra.Command{
    Use:   "switch <branch>",
    Short: "Switch to an existing branch.",
	Long: "Will switch all three odoo repositories to the specified branch or version.",
	Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
    	input, _ := strings.CutPrefix(args[0], "odoo-dev:")
		version := lib.DetectVersion(input)
		lib.Debug.Printf("switchCmd: version '%v' was detected", version)
		var repoBranch sync.Map
		lib.ForEachRepository(func (i int, repoName string, repository *lib.Repository) error {
			branchName := input
			if !repository.BranchExists(branchName) {
				branchName = version
				if !repository.BranchExists(branchName) { branchName = lib.FallbackBranch }
			}
			repoBranch.Store(repoName, branchName)
			return nil
		}, true)

		var spinners sync.Map
		ms := lib.NewMultiSpinner()
		lib.ForEachRepository(func (i int, repoName string, repository *lib.Repository) error {
			branchName, _ := repoBranch.Load(repoName)
			text := fmt.Sprintf("[%s] Switching to '%s'", repository.Color(repoName), branchName)
			spinners.Store(repoName, ms.Add(text))
			return nil
		}, false)
		ms.Start()

		errors := lib.ForEachRepository(func (i int, repoName string, repository *lib.Repository) error {
			branchName, _ := repoBranch.Load(repoName)
			if repoName == ".workspace" {
				curBranch := repository.GetCurrentBranch()
				if curBranch != "master" {
					changes, err := repository.GetStatus()
					if err != nil {
						lib.Debug.Printf("in repository %v: error getting changes, won't auto-commit: %v", repoName, err)
					} else if len(changes) > 0 {
						commitMessage := fmt.Sprintf("odv auto-commit %v\n\nBefore switching to '%s'", time.Now().Format(time.RFC3339), input)
						err = repository.CommitAll(commitMessage)
						if err != nil {
							lib.Debug.Printf("in repository %v: error committing changes, won't auto-commit: %v", repoName, err)
						} else {
							lib.Debug.Printf("in repository %v: auto-committed changes with message '%v'", repoName, commitMessage)
						}
					}
				}
			}

			err := repository.SwitchBranch(branchName.(string))
			ls, _ := spinners.Load(repoName)
			spinner := ls.(*lib.LineSpinner)
			if err != nil {
				ms.Fail(spinner)
			} else {
				ms.Done(spinner)
			}
			return err
		}, true)
		ms.Close()
		lib.PrintRepositoryErrors(errors)
	},
}

func init() {
    rootCmd.AddCommand(switchCmd)
}
