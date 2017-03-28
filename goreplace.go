package main

import (
	"fmt"
	"io/ioutil"
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
	Search          string   `short:"s" long:"regex" description:"Search regexp" value-name:"RE"`
	Replace         string   `short:"r" long:"replace" description:"replace found substrings with RE" value-name:"RE"`
	IgnoreCase      bool     `short:"i" long:"ignore-case" description:"ignore pattern case"`
	PlainText       bool     `short:"p" long:"plain" description:"treat pattern as plain text"`
	Verbose         bool     `short:"v" long:"verbose" description:"verbose mode"`
	ShowVersion     bool     `short:"V" long:"version" description:"show version and exit"`
	ShowHelp        bool     `short:"h" long:"help" description:"show this help message"`
}

var argparser = flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash)

func replaceInFile(file string) {
	var err error

    read, err := ioutil.ReadFile(file)
    if err != nil {
        panic(err)
    }

    content := string(read)

    if opts.PlainText {
        if strings.Contains(content, opts.Search) {
            content = strings.Replace(content, opts.Search, opts.Replace, -1)
            writeContentToFile(file, content)
        }
    } else {
        regex := opts.Search

        if opts.IgnoreCase {
            regex = "(?i:" + regex + ")"
        }

        re := regexp.MustCompile(regex)

		if re.MatchString(content) {
		    content = re.ReplaceAllString(content, opts.Replace)
		    writeContentToFile(file, content)
		}
    }

}

func writeContentToFile(file string, content string) {
    var err error
    err = ioutil.WriteFile(file, []byte(content), 0)
    if err != nil {
        panic(err)
    }
}

func main() {
	args, err := argparser.Parse()
    if err != nil {
        panic(err)
        os.Exit(1)
    }

	if opts.ShowVersion {
		fmt.Printf("goreplace %s\n", Version)
		return
	}

	if (opts.ShowHelp) || (opts.Search == "") || (opts.Replace == "") || (len(args) > 0)  {
		fmt.Println("Usage: golang -s <search regex> -r <replace string> file1 file2 file3")
		argparser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

    for i := range args {
        var file string
        file = args[i]

        if opts.Verbose {
            fmt.Printf(" - checking %s\n", file)
        }
        replaceInFile(file)
    }

    os.Exit(0)
}