package main

import (
	"fmt"
	"getMeMod/server/src"
	"getMeMod/server/store/utils/constants"
	"getMeMod/server/utils/logger"
	"os"
)

func main() {
	// Start the server

	fmt.Println("Creating the log file in the executable directory")

	if err := logger.Initialize(constants.LogsDirName); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	storePath := constants.StoreDirName
	compactedStorePath := constants.CompactedStoreDirName
	if err := src.StartServer(constants.SocketPath, storePath, compactedStorePath); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
