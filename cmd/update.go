package cmd

import (
	"fmt"
	"sync"
	"strings"
	"slices"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/internal"
)

var updateCmd = &cobra.Command{
    Use:   "update",
    Short: "Update current branch.",
	Long: "Will fetch and pull (ff-only) all version branches in all three odoo repositories.",
    Run: func(cmd *cobra.Command, args []string) {
		var branches sync.Map
		var branchRepo sync.Map
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) error {
			if repoName == ".workspace" { return nil }
			curBranch := repository.GetCurrentBranch()
			branchRepo.Store(repoName, curBranch)
			var branchSlice []string
			for _, branch := range repository.GetBranches() {
				if internal.IsVersionBranch(branch) {
					branchSlice = append(branchSlice, branch)
				}
			}
			slices.SortFunc(branchSlice, func(a, b string) int {
				aVersion := internal.GetVersion(a)
				bVersion := internal.GetVersion(b)
				comparison := strings.Compare(aVersion, bVersion)
				if comparison != 0 { return -comparison
				} else { return strings.Compare(a, b) }
			})
			branches.Store(repoName, branchSlice)
			return nil
		}, true)

		var spinners sync.Map
		ms := internal.NewMultiSpinner()
		defer ms.Close()
		internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) error {
			_, ok := branches.Load(repoName)
			if ok {
				text := fmt.Sprintf("[%s] Fetching", repository.Color(repoName))
				spinners.Store(repoName, ms.Add(text))
			}
			return nil
		}, false)
		ms.Start()

		errors := internal.ForEachRepository(func (i int, repoName string, repository *internal.Repository) error {
			branchList, ok := branches.Load(repoName)
			if !ok { return nil }
			s, _ := spinners.Load(repoName)
			spinner := s.(*internal.LineSpinner)
			err := repository.Fetch("origin")
			if err != nil {
				ms.Fail(spinner)
				return fmt.Errorf("fetching remote origin: %v", err)
			}
			for _, branchName := range branchList.([]string) {
				if !internal.IsVersionBranch(branchName) { continue }
				ms.UpdateText(spinner, fmt.Sprintf("[%s] Switching '%s'", repository.Color(repoName), branchName))
				err = repository.SwitchBranch(branchName)
				if err != nil {
					ms.Fail(spinner)
					return fmt.Errorf("switching to branch %v: %v", branchName, err)
				}
				ms.UpdateText(spinner, fmt.Sprintf("[%s] Integrating '%s'", repository.Color(repoName), branchName))
				err = repository.IntegrateChangesFromRemote("origin", branchName)
				if err != nil {
					ms.Fail(spinner)
					return fmt.Errorf("integrating changes from remote origin/%v: %v", branchName, err)
				}
			}
			originalBranch, _ := branchRepo.Load(repoName)
			repository.SwitchBranch(originalBranch.(string))
			ms.UpdateText(spinner, fmt.Sprintf("[%s] Updated", repository.Color(repoName)))
			ms.Done(spinner)
			return nil
		}, true)

		internal.PrintRepositoryErrors(errors)
    },
}

func init() {
    rootCmd.AddCommand(updateCmd)
}
