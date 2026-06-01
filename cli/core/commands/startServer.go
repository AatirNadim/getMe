package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

const (
	ServerExecutablePath = "/opt/getme/getMe-server"
	DataDir              = "/var/lib/getMeStore/dataDir"
	SockDir              = "/tmp/getMeStore/sockDir"
	LogDir               = "/tmp/getMeStore/dumpDir"
	PidFile              = "/tmp/getMeStore/sockDir/server.pid"
)

func checkDirOwnership(dirPath string) error {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist. Please run the setup scripts (e.g. init-data-dir.sh)", dirPath)
	}
	if err != nil {
		return err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("could not get underlying sys stat for %s", dirPath)
	}
	if stat.Uid != uint32(os.Getuid()) {
		return fmt.Errorf("directory %s is not owned by the current user (UID %d). Please ensure you have ownership", dirPath, os.Getuid())
	}
	return nil
}

func isServerRunning() bool {
	pidBytes, err := os.ReadFile(PidFile)
	if err != nil {
		return false // PID file doesn't exist or is unreadable
	}
	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	// Signal 0 checks if the process is alive
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

var StartServerCmd = &cobra.Command{
	Use:   "start-server",
	Short: "Start the getMe server as a background daemon",
	Long:  "Starts the getMe server daemon, checks required directories, and stores the PID.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Check if server binary exists
		execPath := os.Getenv("GETME_SERVER_BIN")
		if execPath == "" {
			execPath = ServerExecutablePath
		}
		if _, err := os.Stat(execPath); os.IsNotExist(err) {
			return fmt.Errorf("server executable not found at %s. Please install the server first (or set GETME_SERVER_BIN)", execPath)
		}

		// 2. Ensure /tmp directories exist
		if err := os.MkdirAll(SockDir, 0700); err != nil {
			return fmt.Errorf("failed to create socket directory: %w", err)
		}
		if err := os.MkdirAll(LogDir, 0700); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		// 3. Verify directory ownership
		dirsToCheck := []string{DataDir, SockDir, LogDir}
		for _, dir := range dirsToCheck {
			if err := checkDirOwnership(dir); err != nil {
				return err
			}
		}

		// 4. Ensure server is not already running
		if isServerRunning() {
			return fmt.Errorf("server is already running. Check %s for the PID", PidFile)
		}

		// Clean up a stale PID file if it exists but the process is dead
		os.Remove(PidFile)

		// 5. Start the server daemon
		execCmd := exec.Command(execPath)
		execCmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

		// Detach standard streams so it doesn't block the CLI
		execCmd.Stdin = nil
		execCmd.Stdout = nil
		execCmd.Stderr = nil

		fmt.Println("Starting getMe server...")
		if err := execCmd.Start(); err != nil {
			return fmt.Errorf("failed to start server daemon: %w", err)
		}

		// 6. Write PID to file
		pidStr := strconv.Itoa(execCmd.Process.Pid)
		if err := os.WriteFile(PidFile, []byte(pidStr), 0600); err != nil {
			return fmt.Errorf("server started (PID %s), but failed to write PID file: %w", pidStr, err)
		}

		fmt.Printf("Server successfully started with PID: %s\n", pidStr)
		return nil
	},
}
