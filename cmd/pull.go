package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
)

var pullCmd = &cobra.Command{
    Use:   "pull",
    Short: "pulls current branch.",
	Long: "Will pull (ff-only) the current branch in all three odoo repositories.",
    Run: func(cmd *cobra.Command, args []string) {
		var branchRepo sync.Map
		lib.ForEachRepository(func (i int, repoName string, repository *lib.Repository) error {
			if repoName == ".workspace" { return nil }
			curBranch := repository.GetCurrentBranch()
			if lib.IsVersionBranch(curBranch) {
				branchRepo.Store(repoName, curBranch)
			} else {
				lib.Debug.Printf("in repository %v: current branch %v is not a version branch, skipping sync", repoName, curBranch)
			}
			return nil
		}, true)

		var spinners sync.Map
		ms := lib.NewMultiSpinner()
		defer ms.Close()
		lib.ForEachRepository(func (i int, repoName string, repository *lib.Repository) error {
			branchName, ok := branchRepo.Load(repoName)
			if ok {
				text := fmt.Sprintf("[%s] Pulling '%s'", repository.Color(repoName), branchName)
				spinners.Store(repoName, ms.Add(text))
			}
			return nil
		}, false)
		ms.Start()

		errors := lib.ForEachRepository(func (i int, repoName string, repository *lib.Repository) error {
			curBranch, ok := branchRepo.Load(repoName)
			if !ok { return nil }
			branchName := curBranch.(string)
			s, _ := spinners.Load(repoName)
			spinner := s.(*lib.LineSpinner)
			err := repository.Pull()
			if err != nil {
				ms.Fail(spinner)
				return fmt.Errorf("pulling branch %v: %v", branchName, err)
			}
			ms.UpdateText(spinner, fmt.Sprintf("[%s] Pulled '%s'", repository.Color(repoName), branchName))
			ms.Done(spinner)
			return nil
		}, true)

		lib.PrintRepositoryErrors(errors)
    },
}

func init() {
    rootCmd.AddCommand(pullCmd)
}
