# go-replace

[![GitHub release](https://img.shields.io/github/release/webdevops/go-replace.svg)](https://github.com/webdevops/go-replace/releases)
[![license](https://img.shields.io/github/license/webdevops/go-replace.svg)](https://github.com/webdevops/go-replace/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/webdevops/go-replace.svg?branch=master)](https://travis-ci.org/webdevops/go-replace)
[![Github All Releases](https://img.shields.io/github/downloads/webdevops/go-replace/total.svg)]()
[![Github Releases](https://img.shields.io/github/downloads/webdevops/go-replace/latest/total.svg)]()

Cli utility for replacing text in files, written in golang and compiled for usage in Docker images

Inspired by https://github.com/piranha/goreplace

## Usage

```
Usage:
  go-replace

Application Options:
  -m, --mode=[replace|line|lineinfile|template] replacement mode - replace: replace match with term; line: replace line with term;
                                                lineinfile: replace line with term or if not found append to term to file; template: parse
                                                content as golang template, search value have to start uppercase (default: replace)
  -s, --search=                                 search term
  -r, --replace=                                replacement term
  -i, --case-insensitive                        ignore case of pattern to match upper and lowercase characters
      --stdin                                   process stdin as input
  -o, --output=                                 write changes to this file (in one file mode)
      --output-strip-ext=                       strip file extension from written files (also available in multi file mode)
      --once=[keep|unique]                      replace search term only one in a file, keep duplicaes (keep, default) or remove them
                                                (unique)
      --regex                                   treat pattern as regex
      --regex-backrefs                          enable backreferences in replace term
      --regex-posix                             parse regex term as POSIX regex
      --path=                                   use files in this path
      --path-pattern=                           file pattern (* for wildcard, only basename of file)
      --path-regex=                             file pattern (regex, full path)
      --ignore-empty                            ignore empty file list, otherwise this will result in an error
  -v, --verbose                                 verbose mode
      --dry-run                                 dry run mode
  -V, --version                                 show version and exit
  -h, --help                                    show this help message
```

Files must be specified as arguments and will be overwritten after parsing. If you want an alternative location for
saving the file the argument can be specified as `source:destination`, eg.
`go-replace -s foobar -r barfoo daemon.conf.tmpl:daemon.conf`.

If `--path` (with or without `--path-pattern` or `--path-regex`) the files inside path are used as source and will
be overwritten. If `daemon.conf.tmpl` should be written as `daemon.conf` the option `--output-strip-ext=.tmpl` will do
this based on the source file name.


| Mode       | Description                                                                                                                                                    |
|:-----------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------|
| replace    | Replace search term inside one line with replacement.                                                                                                          |
| line       | Replace line (if matched term is inside) with replacement.                                                                                                     |
| lineinfile | Replace line (if matched term is inside) with replacement. If no match is found in the whole file the line will be appended to the bottom of the file.         |
| template   | Parse content as [golang template](https://golang.org/pkg/text/template/), arguments are available via `{{.Arg.Name}}` or environment vars via `{{.Env.Name}}` |


### Examples

| Command                                                            | Description                                                                                      |
|:-------------------------------------------------------------------|:-------------------------------------------------------------------------------------------------|
| `go-replace -s foobar -r barfoo file1 file2`                       | Replaces `foobar` to `barfoo` in file1 and file2                                                 |
| `go-replace --regex -s 'foo.*' -r barfoo file1 file2`               | Replaces the regex `foo.*` to `barfoo` in file1 and file2                                        |
| `go-replace --regex --ignore-case -s 'foo.*' -r barfoo file1 file2` | Replaces the regex `foo.*` (and ignore case) to `barfoo` in file1 and file2                      |
| `go-replace --mode=line -s 'foobar' -r barfoo file1 file2`          | Replaces all lines with content `foobar` to `barfoo` (whole line) in file1 and file2             |
| `go-replace -s 'foobar' -r barfoo --path=./ --path-pattern='*.txt'` | Replaces all lines with content `foobar` to `barfoo` (whole line) in *.txt files in current path |

### Example with golang templates

Configuration file `daemon.conf.tmpl`:
```
<VirtualHost ...>
    ServerName {{.Env.SERVERNAME}}
    DocumentRoot {{.Env.DOCUMENTROOT}}
<VirtualHost>

```

Process file with:

```bash
export SERVERNAME=www.foobar.example
export DOCUMENTROOT=/var/www/foobar.example/
go-replace --mode=template daemon.conf.tmpl:daemon.conf
```

Reuslt file `daemon.conf`:
```
<VirtualHost ...>
    ServerName www.foobar.example
    DocumentRoot /var/www/foobar.example/
<VirtualHost>

```

## Installation

```bash
GOREPLACE_VERSION=0.6.0 \
&& wget -O /usr/local/bin/go-replace https://github.com/webdevops/go-replace/releases/download/$GOREPLACE_VERSION/gr-64-linux \
&& chmod +x /usr/local/bin/go-replace
```

