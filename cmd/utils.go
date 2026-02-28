package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var utilsCmd = &cobra.Command{
	Use:   "utils",
	Short: "Utilities for managing odoo.",
}

var utilsKillOdooCmd = &cobra.Command{
	Use:   "kill-odoo",
	Short: "Find and kill the odoo process.",
	Long:  "Finds the pid of the process listening on port 8069 (the default odoo port) and kills it.",
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := exec.Command("lsof", "-ti", ":8069").CombinedOutput()
		if err != nil || len(pid) == 0 {
			cmd.Println("No process found listening on port 8069.")
			return
		}
		output, err := exec.Command("kill", "-9", strings.TrimSpace(string(pid))).CombinedOutput()
		if err != nil {
			cmd.PrintErrln("Failed to kill odoo process:", string(output), err)
			os.Exit(1)
		}
		cmd.Println("Odoo process killed successfully.")
	},
}

var utilsDropdbCmd = &cobra.Command{
	Use:   "dropdbs",
	Short: "Drops all R&D db's in PostgreSQL. Use with caution!",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := exec.Command("psql", "-d", "postgres", "-t", "-c", "SELECT datname FROM pg_database WHERE datname LIKE 'rd-%'").CombinedOutput()
		if err != nil {
			cmd.PrintErrln("Failed to list databases:", err)
			os.Exit(1)
		}
		for dbname := range strings.SplitSeq(string(output), "\n") {
			dbname = strings.TrimSpace(dbname)
			if dbname == "" {
				continue
			}
			err := exec.Command("dropdb", "--if-exists", "--force", dbname).Run()
			if err != nil {
				cmd.PrintErrln("Failed to drop database:", err)
			}
		}
	},
}

func init() {
	utilsCmd.AddCommand(utilsKillOdooCmd)
	utilsCmd.AddCommand(utilsDropdbCmd)

	rootCmd.AddCommand(utilsCmd)
}
