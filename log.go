package main

import (
	"fmt"
	"os"
	"strings"
)

// Log message
func logMessage(message string) {
	if opts.Verbose {
		fmt.Fprint(os.Stderr, message)
	}
}

// Log error object as message
func logError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
}

// Log error object as message
func logFatalErrorAndExit(err error, exitCode int) {
	cmdline := fmt.Sprintf("%s %s", argparser.Command.Name, strings.Join(os.Args[1:], " "))

	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	fmt.Fprintf(os.Stderr, "Command: %s\n", cmdline)

	os.Exit(exitCode)
}
