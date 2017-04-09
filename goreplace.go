package main

import (
    "fmt"
    "sync"
    "errors"
    "bytes"
    "io/ioutil"
    "path/filepath"
    "bufio"
    "os"
    "strings"
    "regexp"
    "text/template"
    flags "github.com/jessevdk/go-flags"
)

const (
    Author  = "webdevops.io"
    Version = "0.5.4"
)

type changeset struct {
    SearchPlain string
    Search      *regexp.Regexp
    Replace     string
    MatchFound  bool
}

type changeresult struct {
    File   fileitem
    Output string
    Status bool
    Error  error
}

type fileitem struct {
    Path        string
}

type templateData struct {
    Arg map[string]string
    Env map[string]string
}

var opts struct {
    Mode                    string   `short:"m"  long:"mode"                          description:"replacement mode - replace: replace match with term; line: replace line with term; lineinfile: replace line with term or if not found append to term to file; template: parse content as golang template, search value have to start uppercase" default:"replace" choice:"replace" choice:"line" choice:"lineinfile" choice:"template"`
    ModeIsReplaceMatch      bool
    ModeIsReplaceLine       bool
    ModeIsLineInFile        bool
    ModeIsTemplate          bool
    Search                  []string `short:"s"  long:"search"                        description:"search term"`
    Replace                 []string `short:"r"  long:"replace"                       description:"replacement term" `
    CaseInsensitive         bool     `short:"i"  long:"case-insensitive"              description:"ignore case of pattern to match upper and lowercase characters"`
    Stdin                   bool     `           long:"stdin"                         description:"process stdin as input"`
    Output                  string   `short:"o"  long:"output"                        description:"write changes to this file (in one file mode)"`
    Once                    string   `           long:"once"                          description:"replace search term only one in a file, keep duplicaes (keep, default) or remove them (unique)" optional:"true" optional-value:"keep" choice:"keep" choice:"unique"`
    Regex                   bool     `           long:"regex"                         description:"treat pattern as regex"`
    RegexBackref            bool     `           long:"regex-backrefs"                description:"enable backreferences in replace term"`
    RegexPosix              bool     `           long:"regex-posix"                   description:"parse regex term as POSIX regex"`
    Path                    string   `           long:"path"                          description:"use files in this path"`
    PathPattern             string   `           long:"path-pattern"                  description:"file pattern (* for wildcard, only basename of file)"`
    PathRegex               string   `           long:"path-regex"                    description:"file pattern (regex, full path)"`
    IgnoreEmpty             bool     `           long:"ignore-empty"                  description:"ignore empty file list, otherwise this will result in an error"`
    Verbose                 bool     `short:"v"  long:"verbose"                       description:"verbose mode"`
    DryRun                  bool     `           long:"dry-run"                       description:"dry run mode"`
    ShowVersion             bool     `short:"V"  long:"version"                       description:"show version and exit"`
    ShowHelp                bool     `short:"h"  long:"help"                          description:"show this help message"`
}

var pathFilterDirectories = []string{"autom4te.cache", "blib", "_build", ".bzr", ".cdv", "cover_db", "CVS", "_darcs", "~.dep", "~.dot", ".git", ".hg", "~.nib", ".pc", "~.plst", "RCS", "SCCS", "_sgbak", ".svn", "_obj", ".idea"}

// Apply changesets to file
func applyChangesetsToFile(fileitem fileitem, changesets []changeset) (string, bool, error) {
    var (
        err error = nil
        output string = ""
        status bool = true
    )

    // try open file
    file, err := os.Open(fileitem.Path)
    if err != nil {
        return output, false, err
    }

    writeBufferToFile := false
    var buffer bytes.Buffer

    r := bufio.NewReader(file)
    line, e := Readln(r)
    for e == nil {
        newLine, lineChanged, skipLine := applyChangesetsToLine(line, changesets)

        if lineChanged || skipLine {
            writeBufferToFile = true
        }

        if !skipLine {
            buffer.WriteString(newLine + "\n")
        }

        line, e = Readln(r)
    }

    // --mode=lineinfile
    if opts.ModeIsLineInFile {
        for i := range changesets {
            changeset := changesets[i]
            if !changeset.MatchFound {
                buffer.WriteString(changeset.Replace + "\n")
                writeBufferToFile = true
            }
        }
    }

    if writeBufferToFile {
        output, status = writeContentToFile(fileitem, buffer)
    } else {
        output = fmt.Sprintf("%s no match", fileitem.Path)
    }

    return output, status, err
}

// Apply changesets to file
func applyTemplateToFile(fileitem fileitem, changesets []changeset) (string, bool, error) {
    var (
        err error = nil
        output string = ""
        status bool = true
    )

    // try open file
    buffer, err := ioutil.ReadFile(fileitem.Path)
    if err != nil {
        return output, false, err
    }

    content := parseContentAsTemplate(string(buffer), changesets)

    output, status = writeContentToFile(fileitem, content)

    return output, status, err
}

func applyChangesetsToLine(line string, changesets []changeset) (string, bool, bool) {
    changed := false
    skipLine := false

    for i := range changesets {
        changeset := changesets[i]

        // --once, only do changeset once if already applied to file
        if opts.Once != "" && changeset.MatchFound {
            // --once=unique, skip matching lines
            if opts.Once == "unique" && searchMatch(line, changeset) {
                // matching line, not writing to buffer as requsted
                skipLine = true
                changed = true
                break
            }
        } else {
            // search and replace
            if searchMatch(line, changeset) {
                // --mode=line or --mode=lineinfile
                if opts.ModeIsReplaceLine || opts.ModeIsLineInFile {
                    // replace whole line with replace term
                    line = changeset.Replace
                } else {
                    // replace only term inside line
                    line = replaceText(line, changeset)
                }

                changesets[i].MatchFound = true
                changed = true
            }
        }
    }

    return line, changed, skipLine
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
func searchMatch(content string, changeset changeset) (bool) {
    if changeset.Search.MatchString(content) {
        return true
    }

    return false
}

// Replace text in whole content based on search options
func replaceText(content string, changeset changeset) (string) {
    // --regex-backrefs
    if opts.RegexBackref {
        return changeset.Search.ReplaceAllString(content, changeset.Replace)
    } else {
        return changeset.Search.ReplaceAllLiteralString(content, changeset.Replace)
    }
}

// Write content to file
func writeContentToFile(fileitem fileitem, content bytes.Buffer) (string, bool) {
    // --dry-run
    if opts.DryRun {
        return content.String(), true
    } else {
        var err error

        filepath := fileitem.Path

        // --output
        if (opts.Output != "") {
            filepath = opts.Output
        }

        err = ioutil.WriteFile(filepath, content.Bytes(), 0644)
        if err != nil {
            panic(err)
        }

        return fmt.Sprintf("%s found and replaced match\n", fileitem.Path), true
    }
}

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

// Build search term
// Compiles regexp if regexp is used
func buildSearchTerm(term string) (*regexp.Regexp) {
    var ret *regexp.Regexp
    var regex string

    // --regex
    if opts.Regex {
        // use search term as regex
        regex = term
    } else {
        // use search term as normal string, escape it for regex usage
        regex = regexp.QuoteMeta(term)
    }

    // --ignore-case
    if opts.CaseInsensitive {
        regex = "(?i:" + regex + ")"
    }

    // --verbose
    if opts.Verbose {
        logMessage(fmt.Sprintf("Using regular expression: %s", regex))
    }

    // --regex-posix
    if opts.RegexPosix {
        ret = regexp.MustCompilePOSIX(regex)
    } else {
        ret = regexp.MustCompile(regex)
    }

    return ret
}

// check if string is contained in an array
func contains(slice []string, item string) bool {
    set := make(map[string]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }

    _, ok := set[item]
    return ok
}

// search files in path
func searchFilesInPath(path string, callback func(os.FileInfo, string)) {
        var pathRegex *regexp.Regexp

        // --path-regex
        if (opts.PathRegex != "") {
            pathRegex = regexp.MustCompile(opts.PathRegex)
        }

        // collect all files
        filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
            filename := f.Name()

            // skip directories
            if f.IsDir() {
                if contains(pathFilterDirectories, f.Name()) {
                    return filepath.SkipDir
                }

                return nil
            }

            // --path-pattern
            if (opts.PathPattern != "") {
                matched, _ := filepath.Match(opts.PathPattern, filename)
                if (!matched) {
                    return nil
                }
            }

            // --path-regex
            if pathRegex != nil {
                if (!pathRegex.MatchString(path)) {
                    return nil
                }
            }

            callback(f, path)
            return nil
        })
}

// handle special cli options
// eg. --help
//     --version
//     --path
//     --mode=...
func handleSpecialCliOptions(args []string) ([]string) {
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

    // --mode
    switch mode := opts.Mode; mode {
        case "replace":
            opts.ModeIsReplaceMatch = true
            opts.ModeIsReplaceLine = false
            opts.ModeIsLineInFile = false
            opts.ModeIsTemplate = false
        case "line":
            opts.ModeIsReplaceMatch = false
            opts.ModeIsReplaceLine = true
            opts.ModeIsLineInFile = false
            opts.ModeIsTemplate = false
        case "lineinfile":
            opts.ModeIsReplaceMatch = false
            opts.ModeIsReplaceLine = false
            opts.ModeIsLineInFile = true
            opts.ModeIsTemplate = false
        case "template":
            opts.ModeIsReplaceMatch = false
            opts.ModeIsReplaceLine = false
            opts.ModeIsLineInFile = false
            opts.ModeIsTemplate = true
    }

    // --path
    if (opts.Path != "") {
        searchFilesInPath(opts.Path, func(f os.FileInfo, path string) {
            args = append(args, path)
        })
    }

    // --output
    if (opts.Output != "" && len(args) > 1) {
        logFatalErrorAndExit(errors.New("Only one file is allowed when using --output"), 1)
    }

    return args
}

func actionProcessStdinReplace(changesets []changeset) (int) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        line := scanner.Text()

        newLine, _, skipLine := applyChangesetsToLine(line, changesets)

        if !skipLine {
            fmt.Println(newLine)
        }
    }

    return 0
}


func generateTemplateData(changesets []changeset) (templateData) {
    // init
    var ret templateData
    ret.Arg = make(map[string]string)
    ret.Env = make(map[string]string)

    // add changesets
    for i := range changesets {
        changeset := changesets[i]
        ret.Arg[changeset.SearchPlain] = changeset.Replace
    }

    // add env variables
    for _, e := range os.Environ() {
        pair := strings.Split(e, "=")

        envKey := pair[0]
        envValue := strings.Join(pair[1:], "=")

        ret.Env[envKey] = envValue
    }

    return ret
}

func parseContentAsTemplate(templateContent string, changesets []changeset) bytes.Buffer {
    var content bytes.Buffer
    data := generateTemplateData(changesets)

    tmpl, err := template.New("template").Parse(templateContent)
    if err != nil {
        logFatalErrorAndExit(err, 1)
    }
    tmpl.Option("missingkey=zero")
    err = tmpl.Execute(&content, &data)
    if err != nil {
        logFatalErrorAndExit(err, 1)
    }

    return content
}

func actionProcessStdinTemplate(changesets []changeset) (int) {
    var buffer bytes.Buffer

    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        buffer.WriteString(scanner.Text() + "\n")
    }

    content := parseContentAsTemplate(buffer.String(), changesets)
    fmt.Print(content.String())

    return 0
}

func actionProcessFiles(changesets []changeset, args []string) (int) {
    // check if there is at least one file to process
    if (len(args) == 0) {
        if (opts.IgnoreEmpty) {
            // no files found, but we should ignore empty filelist
            logMessage("No files found, requsted to ignore this")
            os.Exit(0)
        } else {
            // no files found, print error and exit with error code
            logFatalErrorAndExit(errors.New("No files specified"), 1)
        }
    }

    results := make(chan changeresult)

    var wg sync.WaitGroup

    // process file list
    for i := range args {
        file := fileitem{args[i]}

        wg.Add(1)
        go func(file fileitem, changesets []changeset) {
            var (
                err error = nil
                output string = ""
                status bool = true
            )

            if opts.ModeIsTemplate {
                output, status, err = applyTemplateToFile(file, changesets)
            } else {
                output, status, err = applyChangesetsToFile(file, changesets)
            }
            results <- changeresult{file, output, status, err}
            wg.Done()
        } (file, changesets);
    }

    // wait for all changes to be processed
    go func() {
        wg.Wait()
        close(results)
    }()

    // show results
    errorCount := 0
    for result := range results {
        if result.Error != nil {
            logError(result.Error)
            errorCount++
        } else if opts.Verbose {
            title := fmt.Sprintf("%s:", result.File.Path)

            fmt.Fprintln(os.Stderr, "")
            fmt.Fprintln(os.Stderr, title)
            fmt.Fprintln(os.Stderr, strings.Repeat("-", len(title)))
            fmt.Fprintln(os.Stderr, "")
            fmt.Fprintln(os.Stderr, result.Output)
            fmt.Fprintln(os.Stderr, "")
        }
    }

    if errorCount >= 1 {
        fmt.Fprintln(os.Stderr, fmt.Sprintf("[ERROR] %s failed with %d error(s)", argparser.Command.Name, errorCount))
        return 1
    }

    return 0
}

func buildChangesets() ([]changeset){
    var changesets []changeset

    if !opts.ModeIsTemplate {
        if len(opts.Search) == 0 || len(opts.Replace) == 0 {
            // error: unequal numbers of search and replace options
            logFatalErrorAndExit(errors.New("Missing either --search or --replace for this mode"), 1)
        }
    }

    // check if search and replace options have equal lenght (equal number of options)
    if len(opts.Search) != len(opts.Replace) {
        // error: unequal numbers of search and replace options
        logFatalErrorAndExit(errors.New("Unequal numbers of search or replace options"), 1)
    }

    // build changesets
    for i := range opts.Search {
        search := opts.Search[i]
        replace := opts.Replace[i]

        changeset := changeset{search, buildSearchTerm(search), replace, false}
        changesets = append(changesets, changeset)
    }

    return changesets
}

var argparser *flags.Parser
func main() {
    argparser = flags.NewParser(&opts, flags.PassDoubleDash)
    args, err := argparser.Parse()

    args = handleSpecialCliOptions(args)

    // check if there is an parse error
    if err != nil {
        logFatalErrorAndExit(err, 1)
    }

    changesets := buildChangesets()

    exitMode := 0
    if opts.Stdin {
        if opts.ModeIsTemplate {
            // use stdin as input
            exitMode = actionProcessStdinTemplate(changesets)
        } else {
            // use stdin as input
            exitMode = actionProcessStdinReplace(changesets)
        }
    } else {
        // use and process files (see args)
        exitMode = actionProcessFiles(changesets, args)
    }

    os.Exit(exitMode)
}
