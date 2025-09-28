package logger


import (
	"fmt"
	"os"
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

// core printer
func printMessage(title string, color string, message []any) {
	fmt.Fprintf(os.Stdout, "%s%s[%s]%s %s%s%s\n",
		color, bold, title, reset, color, fmt.Sprint(message), reset,
	)
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
