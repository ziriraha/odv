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
	Use:   "drop <dbname|prefix>",
	Short: "Drops all R&D db's in PostgreSQL. Use with caution!",
	Long:  "Drops databases in PSQL and Filestore. If no args are given with --all, it drops the 'rd-*' databases. If a prefix is given, it drops all databases starting with that prefix.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteAll, _ := cmd.Flags().GetBool("all")
		if deleteAll {
			prefix := "rd-"
			if len(args) == 1 {
				prefix = args[0]
			}
			dbsToDelete, err := lib.ListDBs(prefix)
			if err != nil {
				cmd.PrintErrln("Failed to list databases:", err)
				os.Exit(1)
			}
			for _, dbname := range dbsToDelete {
				err := lib.DropDB(dbname)
				if err != nil {
					cmd.PrintErrf("Failed to drop database %s: %v\n", dbname, err)
				} else {
					cmd.Printf("Dropped database: %s\n", dbname)
				}
			}
			if len(dbsToDelete) == 0 {
				cmd.Printf("No databases found with prefix '%s'.\n", prefix)
			}
		} else {
			if len(args) == 0 {
				cmd.PrintErrln("Please provide a database name or prefix, or use --all to drop all databases.")
				os.Exit(1)
			}
			dbname := args[0]
			err := lib.DropDB(dbname)
			if err != nil {
				cmd.PrintErrf("Failed to drop database %s: %v\n", dbname, err)
			}
		}
	},
}

var dbDuplicateCmd = &cobra.Command{
	Use:   "duplicate <source_db> <new_db>",
	Short: "Duplicates an existing database.",
	Long:  "Creates a new database by duplicating an existing one. The source database must exist, and the new database must not already exist.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		sourceDB := args[0]
		newDB := args[1]
		err := lib.DuplicateDB(sourceDB, newDB)
		if err != nil {
			cmd.PrintErrf("Failed to duplicate database from %s to %s: %v\n", sourceDB, newDB, err)
			os.Exit(1)
		} else {
			cmd.Printf("Successfully duplicated database from %s to %s\n", sourceDB, newDB)
		}
	},
}

func init() {
	dbDropCmd.Flags().BoolP("all", "a", false, "Drop all databases")
	dbCmd.AddCommand(dbDropCmd)

	dbCmd.AddCommand(dbDuplicateCmd)

	rootCmd.AddCommand(dbCmd)
}
