package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var utilsCmd = &cobra.Command{
	Use:   "utils",
	Short: "Utilities for managing odoo.",
}

func findKillOdooProcess() error {
	pid, err := exec.Command("lsof", "-ti", ":8069").CombinedOutput()
	if err != nil || len(pid) == 0 {
		return fmt.Errorf("no process found listening on port 8069")
	}
	pidInt, err := strconv.Atoi(strings.TrimSpace(string(pid)))
	if err != nil {
		return fmt.Errorf("invalid pid: %v", err)
	}
	process, err := os.FindProcess(pidInt)
	if err != nil {
		return fmt.Errorf("could not find process: %v", err)
	}
	err = process.Kill()
	if err != nil {
		return fmt.Errorf("failed to kill process: %v", err)
	}
	return nil
}

var utilsKillOdooCmd = &cobra.Command{
	Use:   "kill-odoo",
	Short: "Find and kill the odoo process.",
	Long:  "Finds the pid of the process listening on port 8069 (the default odoo port) and kills it.",
	Run: func(cmd *cobra.Command, args []string) {
		err := findKillOdooProcess()
		if err != nil {
			cmd.PrintErrln("Failed to kill odoo process:", err)
			os.Exit(1)
		}
		cmd.Println("Odoo process killed successfully.")
	},
}

func init() {
	utilsCmd.AddCommand(utilsKillOdooCmd)

	rootCmd.AddCommand(utilsCmd)
}
