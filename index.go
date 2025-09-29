package main

import (
	"fmt"
	"getMeMod/server"
	"getMeMod/store/utils/constants"
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
	if err := server.StartServer(socketPath, storePath); err != nil {
		fmt.Println("Error starting server:", err)
	}
}