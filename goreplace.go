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
    Version = "0.2.1"
)

var opts struct {
    Search            string   `short:"s"  long:"search"       required:"true"  description:"search term"`
    SearchRegex       *regexp.Regexp
    Replace           string   `short:"r"  long:"replace"      required:"true"  description:"replacement term" `
    IgnoreCase        bool     `short:"i"  long:"ignore-case"                   description:"ignore pattern case"`
    ReplaceLine       bool     `           long:"replace-line"                  description:"replace whole line instead of only match"`
    Regex             bool     `           long:"regex"                         description:"treat pattern as regex"`
    RegexBackref      bool     `           long:"regex-backrefs"                description:"enable backreferences in replace term"`
    RegexPosix        bool     `           long:"regex-posix"                   description:"parse regex term as POSIX regex"`
    Verbose           bool     `short:"v"  long:"verbose"                       description:"verbose mode"`
    DryRun            bool     `           long:"dry-run"                       description:"dry run mode"`
    ShowVersion       bool     `short:"V"  long:"version"                       description:"show version and exit"`
    ShowHelp          bool     `short:"h"  long:"help"                          description:"show this help message"`
}

// Replace line (if match is found) in file
func replaceInFile(filepath string) {
    // try open file
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
            // --replace-line
            if opts.ReplaceLine {
                // replace whole line with replace term
                line = opts.Replace
            } else {
                // replace only term inside line
                line = replaceText(line)
            }

            buffer.WriteString(line + "\n")
            replaceStatus = true
        } else {
            buffer.WriteString(line + "\n")
        }

        line, e = Readln(r)
    }

    if replaceStatus {
        writeContentToFile(filepath, buffer)
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
    if opts.SearchRegex.MatchString(content) {
        return true
    }

    return false
}

// Replace text in whole content based on search options
func replaceText(content string) (string) {
    // --regex-backrefs
    if opts.RegexBackref {
        return opts.SearchRegex.ReplaceAllString(content, opts.Replace)
    } else {
        return opts.SearchRegex.ReplaceAllLiteralString(content, opts.Replace)
    }
}

// Write content to file
func writeContentToFile(filepath string, content bytes.Buffer) {
    // --dry-run
    if opts.DryRun {
        title := fmt.Sprintf("%s:", filepath)

        fmt.Println()
        fmt.Println(title)
        fmt.Println(strings.Repeat("-", len(title)))
        fmt.Println(content.String())
        fmt.Println()
    } else {
        var err error
        err = ioutil.WriteFile(filepath, content.Bytes(), 0)
        if err != nil {
            panic(err)
        }

        logMessage(fmt.Sprintf("%s found and replaced match", filepath))
    }
}


// Log message
func logMessage(message string) {
    if opts.Verbose {
        fmt.Println(message)
    }
}

// Log error object as message
func logError(err error) {
    fmt.Printf("Error: %s\n", err)
}

// Build search term
// Compiles regexp if regexp is used
func buildSearchTerm() {
    var regex string

    // --regex
    if opts.Regex {
        // use search term as regex
        regex = opts.Search
    } else {
        // use search term as normal string, escape it for regex usage
        regex = regexp.QuoteMeta(opts.Search)
    }

    // --ignore-case
    if opts.IgnoreCase {
        regex = "(?i:" + regex + ")"
    }

    if opts.Verbose {
        logMessage(fmt.Sprintf("Using regular expression: %s", regex))
    }

    if opts.RegexPosix {
        opts.SearchRegex = regexp.MustCompilePOSIX(regex)
    } else {
        opts.SearchRegex = regexp.MustCompile(regex)
    }
}

func handleSpecialCliOptions(argparser *flags.Parser, args []string) {
    // --version
    if (opts.ShowVersion) {
        fmt.Printf("goreplace version %s\n", Version)
        os.Exit(0)
    }

    // --help
    if (opts.ShowHelp) {
        argparser.WriteHelp(os.Stdout)
        os.Exit(1)
    }

    // missing any files
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

    handleSpecialCliOptions(argparser, args)

    if err != nil {
        logError(err)
        fmt.Println()
        argparser.WriteHelp(os.Stdout)
        os.Exit(1)
    }

    buildSearchTerm()

    for i := range args {
        var file string
        file = args[i]

        replaceInFile(file)
    }

    os.Exit(0)
}
