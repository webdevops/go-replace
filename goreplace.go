package main

import (
	"fmt"
	"errors"
	"bytes"
	"io/ioutil"
	"bufio"
	"os"
	"strings"
	"regexp"
	flags "github.com/jessevdk/go-flags"
)

const (
	Author  = "webdevops.io"
	Version = "1.0"
)

var opts struct {
	Search            string   `short:"s"  long:"search"       required:"true"  description:"search term"`
	SearchRegex       *regexp.Regexp
	Replace           string   `short:"r"  long:"replace"      required:"true"  description:"replacement term" `
	IgnoreCase        bool     `short:"i"  long:"ignore-case"                   description:"ignore pattern case"`
	WholeLine         bool     `           long:"whole-line"                    description:"replace whole line"`
	Regex             bool     `           long:"regex"                         description:"treat pattern as regex"`
	Verbose           bool     `short:"v"  long:"verbose"                       description:"verbose mode"`
	DryRun            bool     `           long:"dry-run"                       description:"dry run mode"`
	ShowVersion       bool     `short:"V"  long:"version"                       description:"show version and exit"`
	ShowHelp          bool     `short:"h"  long:"help"                          description:"show this help message"`
}

// Replace content parts in file
func replaceContentInFile(filepath string) {
	var err error

    read, err := ioutil.ReadFile(filepath)
    if err != nil {
        panic(err)
    }

    content, replaceStatus := replaceText(string(read))

    if replaceStatus {
        writeContentToFile(filepath, content)
    } else {
        logMessage(fmt.Sprintf("%s no match", filepath))
    }

}

// Replace line (if match is found) in file
func replaceLineInFile(filepath string) {
    file, err := os.Open(filepath)
    if err != nil {
        panic(err)
    }

    replaceStatus := false
    var buffer bytes.Buffer

    r := bufio.NewReader(file)
    line, e := Readln(r)
    for e == nil {
        if searchMatch(line) {
            buffer.WriteString(opts.Replace + "\n")
            replaceStatus = true
        } else {
            buffer.WriteString(line + "\n")
        }

        line, e = Readln(r)
    }

    if replaceStatus {
        writeContentToFile(filepath, buffer.String())
    } else {
        logMessage(fmt.Sprintf("%s no match", filepath))
    }
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (string, error) {
  var (isPrefix bool = true
       err error = nil
       line, ln []byte
      )
  for isPrefix && err == nil {
      line, isPrefix, err = r.ReadLine()
      ln = append(ln, line...)
  }
  return string(ln),err
}


// Checks if there is a match in content, based on search options
func searchMatch(content string) (bool) {
    if opts.Regex {
        if opts.SearchRegex.MatchString(content) {
            return true
        }
    } else {
        if strings.Contains(content, opts.Search) {
            return true
        }
    }

    return false
}

// Replace text in whole content based on search options
func replaceText(content string) (string, bool) {
    status := false

    if searchMatch(content) {
        status = true

        if opts.Regex {
            content = opts.SearchRegex.ReplaceAllString(content, opts.Replace)
        } else {
            content = strings.Replace(content, opts.Search, opts.Replace, -1)
        }
    }

    return content, status
}

// Write content to file
func writeContentToFile(filepath string, content string) {
    if opts.DryRun {
        title := fmt.Sprintf("%s:", filepath)

        fmt.Println(title)
        fmt.Println(strings.Repeat("-", len(title)))
        fmt.Println(content)
        fmt.Println()
        fmt.Println()
    } else {
        var err error
        err = ioutil.WriteFile(filepath, []byte(content), 0)
        if err != nil {
            panic(err)
        }
    }
}


// Log message
func logMessage(message string) {
    if opts.Verbose {
        fmt.Println(message)
    }
}

func logError(err error) {
    fmt.Printf("Error: %s\n", err)
}

// Process search option
// Compiles regexp if regexp is used
func processSearch() {
    if opts.Regex {
        regex := opts.Search

        if opts.IgnoreCase {
            regex = "(?i:" + regex + ")"
        }

        opts.SearchRegex = regexp.MustCompile(regex)
    }
}

func handleSpecialOptions(argparser *flags.Parser, args []string) {
    if (opts.ShowVersion) {
        fmt.Printf("goreplace %s\n", Version)
        os.Exit(0)
    }

    if (opts.ShowHelp) {
		argparser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	if (len(args) == 0) {
	    err := errors.New("No files specified")
	    logError(err)
	    fmt.Println()
		argparser.WriteHelp(os.Stdout)
		os.Exit(1)
	}
}

func main() {
    var argparser = flags.NewParser(&opts, flags.PassDoubleDash)
	args, err := argparser.Parse()

    if err != nil {
        handleSpecialOptions(argparser, args)

        logError(err)
        fmt.Println()
        argparser.WriteHelp(os.Stdout)
        os.Exit(1)
    }

	handleSpecialOptions(argparser, args)

	processSearch()

    for i := range args {
        var file string
        file = args[i]

        if opts.WholeLine {
            replaceLineInFile(file)
        } else {
            replaceContentInFile(file)
        }
    }

    os.Exit(0)
}