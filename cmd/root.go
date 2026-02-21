package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
)

var rootCmd = &cobra.Command{
	Use:   "odv",
	Short: "An all in one tool for Odoo development.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		lib.InitializeConfiguration()
		lib.PrefetchAllBranches()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
