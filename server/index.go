package main

import (
	"flag"
	"fmt"

	"github.com/AatirNadim/getMe/server/src"
	"github.com/AatirNadim/getMe/server/store/utils/constants"
)

func main() {
	// check whether the logging is disabled via command line flag
	loggingDisabled := flag.Bool("logging_disabled", false, "disable persistent logging output")
	flag.Parse()

	if err := src.StartServer(constants.SocketPath, constants.StoreDirName, constants.CompactedStoreDirName, loggingDisabled); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
