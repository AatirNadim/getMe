package main

import (
	"flag"
	"fmt"
	"getMeMod/server/src"
	"getMeMod/server/store/utils/constants"
)

func main() {
	// check whether the logging is disabled via command line flag
	loggingDisabled := flag.Bool("logging_disabled", false, "disable persistent logging output")
	flag.Parse()

	storePath := constants.StoreDirName
	compactedStorePath := constants.CompactedStoreDirName
	if err := src.StartServer(constants.SocketPath, storePath, compactedStorePath, loggingDisabled); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
