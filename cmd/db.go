package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/ziriraha/odv/lib"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management for Odoo.",
}

var dbDropCmd = &cobra.Command{
	Use:   "dropall",
	Short: "Drops all Odoo db's. Use with caution!",
	Long:  "Drops odoo databases using odoo's db drop utility.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prefix := lib.GetDefaultDBPrefix()
		if len(args) == 1 {
			prefix = args[0]
		}
		dbsToDelete, err := lib.ListDBs(prefix)
		if err != nil {
			cmd.PrintErrln("Failed to list databases:", err)
			os.Exit(1)
		}
		if len(dbsToDelete) == 0 {
			cmd.Printf("No databases found with prefix '%s'.\n", prefix)
		} else {
			err := lib.DropDBs(dbsToDelete...)
			if err != nil {
				cmd.PrintErrf("Failed to drop databases %s: %v\n", dbsToDelete, err)
			} else {
				cmd.Printf("Dropped databases: %s\n", dbsToDelete)
			}
		}
	},
}

var dbListCmd = &cobra.Command{
	Use:   "list [prefix]",
	Short: "Lists all R&D databases in PostgreSQL.",
	Long:  "Lists all databases in PostgreSQL that start with 'rd-'.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prefix := lib.GetDefaultDBPrefix()
		if len(args) == 1 {
			prefix = args[0]
		}
		dbs, err := lib.ListDBs(prefix)
		if err != nil {
			cmd.PrintErrln("Failed to list databases:", err)
			os.Exit(1)
		}
		if len(dbs) == 0 {
			cmd.Printf("No databases found with the '%s' prefix.\n", prefix)
			return
		}
		for _, db := range dbs {
			cmd.Printf("%s\n", db)
		}
	},
}

func init() {
	dbDropCmd.Flags().BoolP("all", "a", false, "Drop all databases")
	dbCmd.AddCommand(dbDropCmd)

	dbCmd.AddCommand(dbListCmd)

	rootCmd.AddCommand(dbCmd)
}
