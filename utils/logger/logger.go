package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// ANSI color codes
const (
	reset = "\033[0m"
	bold  = "\033[1m"

	offwhite = "\033[37m"

	red    = "\033[31m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	green  = "\033[32m"
	cyan   = "\033[36m"
	grey   = "\033[90m"
)

var (
	logFile     *os.File
	logFileOnce sync.Once
	logFileErr  error
)

// getLogFile returns the singleton log file instance
func getLogFile() (*os.File, error) {
	logFileOnce.Do(func() {
		// Create the dump directory if it doesn't exist
		dumpDir := "dump"
		if err := os.MkdirAll(dumpDir, 0755); err != nil {
			logFileErr = fmt.Errorf("failed to create dump directory: %w", err)
			return
		}

		// Open the log file in append mode
		logPath := filepath.Join(dumpDir, "index.log")
		logFile, logFileErr = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if logFileErr != nil {
			logFileErr = fmt.Errorf("failed to open log file: %w", logFileErr)
		}
	})
	return logFile, logFileErr
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
	// Print colored output to stdout for human readability
	fmt.Fprintf(os.Stdout, "%s%s[%s]%s %s%s%s\n",
		color, bold, title, reset, color, fmt.Sprint(message...), reset,
	)

	// Get the singleton log file
	file, err := getLogFile()
	if err != nil {
		// If we can't open the log file, write error to stderr
		fmt.Fprintf(os.Stderr, "level=ERROR msg=%q\n", fmt.Sprintf("Logger error: %v", err))
		return
	}

	// Write in logfmt format for log aggregation (without colors)
	fmt.Fprintf(file, "level=%s msg=%q\n", title, fmt.Sprint(message...))
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
