package main

import (
	"fmt"
	"getMeMod/server/src"
	"getMeMod/server/store/utils/constants"
	// "getMeMod/server/utils/logger"
	"os"
	"path/filepath"
)

func main() {
	// Start the server
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	fmt.Println("Creating the log file in the executable directory")

	// logPath := filepath.Join("/tmp", "getMeStore", "dump", "index.log")
	// if err := logger.Initialize(logPath); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer logger.Close()

	storePath := filepath.Join(homeDir, constants.StoreDirName)
	compactedStorePath := filepath.Join(homeDir, constants.CompactedStoreDirName)
	if err := src.StartServer(constants.SocketPath, storePath, compactedStorePath); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
