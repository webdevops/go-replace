package main

import (
    "fmt"
    "os"
)

// Log message
func logMessage(message string) {
    if opts.Verbose {
        fmt.Fprintln(os.Stderr, message)
    }
}

// Log error object as message
func logError(err error) {
    fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: %s\n", err))
}


// Log error object as message
func logFatalErrorAndExit(err error, exitCode int) {
    fmt.Fprintln(os.Stderr, fmt.Sprintf("Error: %s\n", err))
    os.Exit(exitCode)
}
