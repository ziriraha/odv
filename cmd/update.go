package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/internal"
)

var updateCmd = &cobra.Command{
    Use:   "update",
    Short: "Update current branch.",
	Long: "Will fetch and pull (ff-only) the current branch in all three odoo repositories.",
    Run: func(cmd *cobra.Command, args []string) {
		var branchRepo sync.Map
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			if repoName == ".vscode" { return }
			curBranch := repository.GetCurrentBranch()
			if !internal.IsVersionBranch(curBranch) {
				internal.Debug.Printf("in repository %v: current branch %v is not a version branch, skipping update", repoName, curBranch)
				return
			}
			branchRepo.Store(repoName, curBranch)
		}, true)

		var spinners sync.Map
		ms := internal.NewMultiSpinner()
		defer ms.Close()
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			branchName, ok := branchRepo.Load(repoName)
			if !ok { return }
			text := fmt.Sprintf("[%s] Fetching '%s'", repository.Color(repoName), branchName)
			spinners.Store(repoName, ms.Add(text))
		}, false)
		ms.Start()

		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) {
			curBranch, ok := branchRepo.Load(repoName)
			if !ok { return }
			branchName := curBranch.(string)
			s, _ := spinners.Load(repoName)
			spinner := s.(*internal.LineSpinner)
			if !ok { return }
			err := repository.Fetch(branchName)
			if err != nil {
				ms.AddOnClose(func() {
					internal.Error.Printf("in repository %v: fetching branch %v: %v", repoName, branchName, err)
				})
				ms.Fail(spinner)
				return
			}
			newSpinnerText := fmt.Sprintf("[%s] Pulling '%s'", repository.Color(repoName), branchName)
			ms.UpdateText(spinner, newSpinnerText)
			err = repository.Pull()
			if err != nil {
				ms.AddOnClose(func() {
					internal.Error.Printf("in repository %v: pulling branch %v: %v", repoName, branchName, err)
				})
				ms.Fail(spinner)
			} else {
				newSpinnerText := fmt.Sprintf("[%s] Updating '%s'", repository.Color(repoName), branchName)
				ms.UpdateText(spinner, newSpinnerText)
				ms.Done(spinner)
			}
		}, true)
	},
}

func init() {
    rootCmd.AddCommand(updateCmd)
}
