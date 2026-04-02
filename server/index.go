package main

import (
	"flag"
	"fmt"

	"github.com/AatirNadim/getMe/server/src"
	"github.com/AatirNadim/getMe/server/store/utils/constants"
)

func main() {
	// check whether the logging is disabled via command line flag


	// use the build-in flag `-h` or `--help` to display usage information about the logging flags

	var loggingDisabled *bool
	var logToStdout *bool

	loggingDisabled = flag.Bool("logging_disabled", false, "disable persistent logging output")
	loggingDisabledShort := flag.Bool("d", false, "disable persistent logging output (shorthand)")

	// if the short flag is set, it takes precedence over the long flag
	if *loggingDisabledShort {
		loggingDisabled = loggingDisabledShort
	}

	logToStdout = flag.Bool("log_to_stdout", false, "direct logs to stdout instead of a file")
	logToStdoutShort := flag.Bool("s", false, "direct logs to stdout instead of a file (shorthand)")

	// if the short flag is set, it takes precedence over the long flag
	if *logToStdoutShort {
		logToStdout = logToStdoutShort
	}

	flag.Parse()

	if err := src.StartServer(constants.SocketPath, constants.StoreDirName, constants.CompactedStoreDirName, loggingDisabled, logToStdout); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
