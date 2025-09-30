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
	socketPath := constants.SocketPath
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	storePath := filepath.Join(homeDir, constants.StoreDirName)
	if err := src.StartServer(socketPath, storePath); err != nil {
		fmt.Println("Error starting server:", err)
	}
}