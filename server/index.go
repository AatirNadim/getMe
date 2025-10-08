package main

import (
	"fmt"
	"getMeMod/server/src"
	"getMeMod/server/store/utils/constants"
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
	storePath := filepath.Join(homeDir, constants.StoreDirName)
	compactedStorePath := filepath.Join(homeDir, constants.CompactedStoreDirName)
	if err := src.StartServer(constants.SocketPath, storePath, compactedStorePath); err != nil {
		fmt.Println("Error starting server:", err)
	}
}