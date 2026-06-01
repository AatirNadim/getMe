package commands

import (
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

var StopServerCmd = &cobra.Command{
	Use:   "stop-server",
	Short: "Stop the background getMe server daemon",
	Long:  "Stops the getMe server daemon using its PID file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Check if PID file exists
		pidBytes, err := os.ReadFile(PidFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Server is not currently running (no PID file found).")
				return nil
			}
			return fmt.Errorf("failed to read PID file: %w", err)
		}

		pid, err := strconv.Atoi(string(pidBytes))
		if err != nil {
			return fmt.Errorf("invalid PID in pidfile: %w", err)
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			return fmt.Errorf("failed to find process with PID %d: %w", pid, err)
		}

		// Check if it's actually alive
		if err := process.Signal(syscall.Signal(0)); err != nil {
			fmt.Printf("Process %d does not exist or is inaccessible. Cleaning up stale PID file.\n", pid)
			_ = os.Remove(PidFile)
			return nil
		}

		// 2. Send SIGTERM
		fmt.Printf("Stopping server (PID %d)...\n", pid)
		if err := process.Signal(syscall.SIGTERM); err != nil {
			return fmt.Errorf("failed to send SIGTERM to process: %w", err)
		}

		// 3. Remove PID file
		if err := os.Remove(PidFile); err != nil {
			fmt.Printf("Warning: Failed to remove PID file %s: %v\n", PidFile, err)
		}

		fmt.Println("Server successfully stopped.")
		return nil
	},
}
