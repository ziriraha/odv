package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odoodev/internal"
)

var rootCmd = &cobra.Command{
	Use:   "odoodev",
	Short: "An all in one tool for Odoo development.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
        internal.SetupLoggers(debug)
		internal.InitializeConfiguration()
    },
}

func init() {
    rootCmd.PersistentFlags().Bool("debug", false, "Enable debug output")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil { os.Exit(1) }
}
