package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/remeh/sizedwaitgroup"
)

const (
	Author = "webdevops.io"
)

var (
	// Git version information
	gitCommit = "<unknown>"
	gitTag    = "<unknown>"
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
	Path   string
	Output string
}

var opts struct {
	ThreadCount        int    `           long:"threads"                       description:"Set thread concurrency for replacing in multiple files at same time" default:"20"`
	Mode               string `short:"m"  long:"mode"                          description:"replacement mode - replace: replace match with term; line: replace line with term; lineinfile: replace line with term or if not found append to term to file; template: parse content as golang template, search value have to start uppercase" default:"replace" choice:"replace" choice:"line" choice:"lineinfile" choice:"template"`
	ModeIsReplaceMatch bool
	ModeIsReplaceLine  bool
	ModeIsLineInFile   bool
	ModeIsTemplate     bool
	Search             []string `short:"s"  long:"search"                        description:"search term"`
	Replace            []string `short:"r"  long:"replace"                       description:"replacement term"`
	LineinfileBefore   string   `           long:"lineinfile-before"             description:"add line before this regex"`
	LineinfileAfter    string   `           long:"lineinfile-after"              description:"add line after this regex"`
	CaseInsensitive    bool     `short:"i"  long:"case-insensitive"              description:"ignore case of pattern to match upper and lowercase characters"`
	Stdin              bool     `           long:"stdin"                         description:"process stdin as input"`
	Output             string   `short:"o"  long:"output"                        description:"write changes to this file (in one file mode)"`
	OutputStripFileExt string   `           long:"output-strip-ext"              description:"strip file extension from written files (also available in multi file mode)"`
	Once               string   `           long:"once"                          description:"replace search term only one in a file, keep duplicaes (keep, default) or remove them (unique)" optional:"true" optional-value:"keep" choice:"keep" choice:"unique"`
	Regex              bool     `           long:"regex"                         description:"treat pattern as regex"`
	RegexBackref       bool     `           long:"regex-backrefs"                description:"enable backreferences in replace term"`
	RegexPosix         bool     `           long:"regex-posix"                   description:"parse regex term as POSIX regex"`
	Path               string   `           long:"path"                          description:"use files in this path"`
	PathPattern        string   `           long:"path-pattern"                  description:"file pattern (* for wildcard, only basename of file)"`
	PathRegex          string   `           long:"path-regex"                    description:"file pattern (regex, full path)"`
	IgnoreEmpty        bool     `           long:"ignore-empty"                  description:"ignore empty file list, otherwise this will result in an error"`
	Verbose            bool     `short:"v"  long:"verbose"                       description:"verbose mode"`
	DryRun             bool     `           long:"dry-run"                       description:"dry run mode"`
	ShowVersion        bool     `short:"V"  long:"version"                       description:"show version and exit"`
	ShowOnlyVersion    bool     `           long:"dumpversion"                   description:"show only version number and exit"`
	ShowHelp           bool     `short:"h"  long:"help"                          description:"show this help message"`
}

var pathFilterDirectories = []string{"autom4te.cache", "blib", "_build", ".bzr", ".cdv", "cover_db", "CVS", "_darcs", "~.dep", "~.dot", ".git", ".hg", "~.nib", ".pc", "~.plst", "RCS", "SCCS", "_sgbak", ".svn", "_obj", ".idea"}

// Apply changesets to file
func applyChangesetsToFile(fileitem fileitem, changesets []changeset) (string, bool, error) {
	var (
		output string = ""
		status bool   = true
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
	file.Close()

	// --mode=lineinfile
	if opts.ModeIsLineInFile {
		lifBuffer, lifStatus := handleLineInFile(changesets, buffer)
		if lifStatus {
			buffer.Reset()
			buffer.WriteString(lifBuffer.String())
			writeBufferToFile = lifStatus
		}
	}

	// --output
	// --output-strip-ext
	// enforcing writing of file (creating new file)
	if opts.Output != "" || opts.OutputStripFileExt != "" {
		writeBufferToFile = true
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
	// try open file
	buffer, err := os.ReadFile(fileitem.Path)
	if err != nil {
		return "", false, err
	}

	content := parseContentAsTemplate(string(buffer), changesets)

	output, status := writeContentToFile(fileitem, content)

	return output, status, err
}

func applyChangesetsToLine(line string, changesets []changeset) (string, bool, bool) {
	changed := false
	skipLine := false

	for i, changeset := range changesets {
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
					if opts.RegexBackref {
						// get match
						line = string(changeset.Search.Find([]byte(line)))

						// replace regex backrefs in match
						line = changeset.Search.ReplaceAllString(line, changeset.Replace)
					} else {
						// replace whole line with replace term
						line = changeset.Replace
					}
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

// Build search term
// Compiles regexp if regexp is used
func buildSearchTerm(term string) *regexp.Regexp {
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

// handle special cli options
// eg. --help
//
//	--version
//	--path
//	--mode=...
func handleSpecialCliOptions(args []string) {
	// --dumpversion
	if opts.ShowOnlyVersion {
		fmt.Println(gitTag)
		os.Exit(0)
	}

	// --version
	if opts.ShowVersion {
		fmt.Printf("go-replace version %s (%s)\n", gitTag, gitCommit)
		fmt.Printf("Copyright (C) 2022 %s\n", Author)
		os.Exit(0)
	}

	// --help
	if opts.ShowHelp {
		argparser.WriteHelp(os.Stdout)
		os.Exit(0)
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

	// --output
	if opts.Output != "" && len(args) > 1 {
		logFatalErrorAndExit(errors.New("Only one file is allowed when using --output"), 1)
	}

	if opts.LineinfileBefore != "" || opts.LineinfileAfter != "" {
		if !opts.ModeIsLineInFile {
			logFatalErrorAndExit(errors.New("--lineinfile-after and --lineinfile-before only valid in --mode=lineinfile"), 1)
		}

		if opts.LineinfileBefore != "" && opts.LineinfileAfter != "" {
			logFatalErrorAndExit(errors.New("Only --lineinfile-after or --lineinfile-before is allowed in --mode=lineinfile"), 1)
		}
	}
}

func actionProcessStdinReplace(changesets []changeset) int {
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

func actionProcessStdinTemplate(changesets []changeset) int {
	var buffer bytes.Buffer

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		buffer.WriteString(scanner.Text() + "\n")
	}

	content := parseContentAsTemplate(buffer.String(), changesets)
	fmt.Print(content.String())

	return 0
}

func actionProcessFiles(changesets []changeset, fileitems []fileitem) int {
	// check if there is at least one file to process
	if len(fileitems) == 0 {
		if opts.IgnoreEmpty {
			// no files found, but we should ignore empty filelist
			logMessage("No files found, requsted to ignore this")
			os.Exit(0)
		} else {
			// no files found, print error and exit with error code
			logFatalErrorAndExit(errors.New("No files specified"), 1)
		}
	}

	swg := sizedwaitgroup.New(8)
	results := make(chan changeresult, len(fileitems))

	// process file list
	for _, file := range fileitems {
		swg.Add()
		go func(file fileitem, changesets []changeset) {
			var (
				err    error
				output string
				status bool
			)

			if opts.ModeIsTemplate {
				output, status, err = applyTemplateToFile(file, changesets)
			} else {
				output, status, err = applyChangesetsToFile(file, changesets)
			}

			results <- changeresult{file, output, status, err}
			swg.Done()
		}(file, changesets)
	}

	// wait for all changes to be processed
	swg.Wait()
	close(results)

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
		fmt.Fprintf(os.Stderr, "[ERROR] %s failed with %d error(s)\n", argparser.Command.Name, errorCount)
		return 1
	}

	return 0
}

func buildChangesets() []changeset {
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

func buildFileitems(args []string) []fileitem {
	var (
		fileitems []fileitem
		file      fileitem
	)

	// Build filelist from arguments
	for _, filepath := range args {
		file = fileitem{filepath, filepath}

		if opts.Output != "" {
			// use specific output
			file.Output = opts.Output
		} else if opts.OutputStripFileExt != "" {
			// remove file ext from saving destination
			file.Output = strings.TrimSuffix(file.Output, opts.OutputStripFileExt)
		} else if strings.Contains(filepath, ":") {
			// argument like "source:destination"
			split := strings.SplitN(filepath, ":", 2)

			file.Path = split[0]
			file.Output = split[1]
		}

		fileitems = append(fileitems, file)
	}

	// --path parsing
	if opts.Path != "" {
		searchFilesInPath(opts.Path, func(f os.FileInfo, filepath string) {
			file := fileitem{filepath, filepath}

			if opts.OutputStripFileExt != "" {
				// remove file ext from saving destination
				file.Output = strings.TrimSuffix(file.Output, opts.OutputStripFileExt)
			}

			// no colon parsing here

			fileitems = append(fileitems, file)
		})
	}

	return fileitems
}

var argparser *flags.Parser

func main() {
	argparser = flags.NewParser(&opts, flags.PassDoubleDash)
	args, err := argparser.Parse()

	handleSpecialCliOptions(args)

	// check if there is an parse error
	if err != nil {
		logFatalErrorAndExit(err, 1)
	}

	changesets := buildChangesets()
	fileitems := buildFileitems(args)

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
		exitMode = actionProcessFiles(changesets, fileitems)
	}

	os.Exit(exitMode)
}
