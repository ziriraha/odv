package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/views"
)

var rootCmd = &cobra.Command{
	Use:   "odv",
	Short: "An all in one tool for Odoo development.",
}

func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.SetErrPrefix(views.ErrorStyle.Render("ERROR "))

	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
