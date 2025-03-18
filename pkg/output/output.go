// Package output provides simple console output utilities.
// It offers functions for displaying informational messages,
// error messages, and fatal errors that terminate the program.
package output

import (
	"fmt"
	"os"
)

// Info prints an informational message to stdout
func Info(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Error prints an error message to stderr
func Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
}

// Fatal prints an error message to stderr and exits with status code 1
func Fatal(format string, args ...interface{}) {
	Error(format, args...)
	os.Exit(1)
}
