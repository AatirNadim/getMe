package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

// ANSI color codes
const (
	offwhite = "\033[37m"

	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
	cyan   = "\033[36m"
)

var (
	logFile         *os.File
	logFilePath     string
	logFileOnce     sync.Once
	logFileErr      error
	loggingDisabled atomic.Bool
	logToStdout     atomic.Bool
)

func Initialize(logPath string, loggingToShell *bool) error {
	logFilePath = logPath
	// updating the logToStdout atomic value based on the provided flag
	if loggingDisabled.Load() {
		fmt.Println("logging is disabled")
		return nil
	}

	if *loggingToShell {
		fmt.Println("Logging to stdout is enabled")
		logToStdout.Store(true)
	} else {
		fmt.Println("Logging to stdout is disabled; logs will be written to file")
		logToStdout.Store(false)
		_, err := getLogFile()
		return err
	}
	// Force initialization

	return nil
}

func Disable() {
	loggingDisabled.Store(true)
}

func Enable() {
	loggingDisabled.Store(false)
}

func DisableLoggingToStdout() {
	logToStdout.Store(false)
}

func EnableLoggingToStdout() {
	logToStdout.Store(true)
}

// getLogFile returns the singleton log file instance
func getLogFile() (*os.File, error) {
	if loggingDisabled.Load() {
		return nil, nil
	}
	logFileOnce.Do(func() {
		if loggingDisabled.Load() {
			fmt.Println("logging is disabled inside logfileonce")
			return
		}
		if logFilePath == "" {
			fmt.Println("log file path is empty")
			logFileErr = fmt.Errorf("logger not initialized: call logger.Initialize() first")
			return
		}

		// Create the directory if it doesn't exist
		dumpDir := filepath.Dir(logFilePath)
		fmt.Println("Creating the dump directory at path:", dumpDir)
		if err := os.MkdirAll(dumpDir, 0755); err != nil {
			logFileErr = fmt.Errorf("failed to create dump directory: %w", err)
			return
		}

		// Open the log file in append mode
		logFile, logFileErr = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if logFileErr != nil {
			logFileErr = fmt.Errorf("failed to open log file: %w", logFileErr)
		}
	})
	if loggingDisabled.Load() {
		return nil, nil
	}
	return logFile, logFileErr
}

func getOutputWriter() (io.Writer, error) {
	if logToStdout.Load() {
		return os.Stdout, nil
	}
	return getLogFile()
}

// Close closes the log file (call this on application shutdown)
func Close() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

// core printer - outputs in logfmt format for Loki/Alloy parsing
func printMessage(title string, color string, message []any) {
	if loggingDisabled.Load() {
		return
	}
	// Print colored output to stdout for human readability
	// fmt.Fprintf(os.Stdout, "%s%s[%s]%s %s%s%s\n",
	// 	color, bold, title, reset, color, fmt.Sprint(message...), reset,
	// )

	// Get the singleton log file

	writer, err := getOutputWriter()
	if err != nil {
		// If we can't open the log file, write error to stderr
		fmt.Fprintf(os.Stderr, "level=ERROR msg=%q\n", fmt.Sprintf("Logger error: %v", err))
		return
	}

	if writer == nil {
		fmt.Println("Logging is disabled, not writing to log file")
		return
	}

	fmt.Fprintf(writer, "level=%s timeStamp=%s msg=%q\n", title, time.Now().Format(time.RFC3339), fmt.Sprint(message...))
}

// Public functions

func Info(msg ...any) {
	printMessage("INFO", offwhite, msg)
}

func Warn(msg ...any) {
	printMessage("WARN", yellow, msg)
}

func Error(msg ...any) {
	printMessage("ERROR", red, msg)
}

func Success(msg ...any) {
	printMessage("SUCCESS", green, msg)
}

func Debug(msg ...any) {
	printMessage("DEBUG", cyan, msg)
}

// func Dim(msg ...any) {
// 	fmt.Fprintf(os.Stdout, "%s%s%s\n", grey, fmt.Sprint(msg), reset)
// }
