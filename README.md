# goreplace

[![GitHub release](https://img.shields.io/github/release/webdevops/goreplace.svg)](https://github.com/webdevops/goreplace/releases)
[![license](https://img.shields.io/github/license/webdevops/goreplace.svg)](https://github.com/webdevops/goreplace/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/webdevops/goreplace.svg?branch=master)](https://travis-ci.org/webdevops/goreplace)
[![Github All Releases](https://img.shields.io/github/downloads/webdevops/goreplace/total.svg)]()
[![Github Releases](https://img.shields.io/github/downloads/webdevops/goreplace/latest/total.svg)]()

Cli utility for replacing text in files, written in golang and compiled for usage in Docker images

Inspired by https://github.com/piranha/goreplace

## Usage

```
Usage:
  goreplace

Application Options:
  -m, --mode=[replace|line|lineinfile|template] replacement mode - replace: replace match with term; line: replace line with term; lineinfile: replace line with term or if not found append to term to file; template:
                                                parse content as golang template, search value have to start uppercase (default: replace)
  -s, --search=                                 search term
  -r, --replace=                                replacement term
  -i, --case-insensitive                        ignore case of pattern to match upper and lowercase characters
      --stdin                                   process stdin as input
  -o, --output=                                 write changes to this file (in one file mode)
      --once=[keep|unique]                      replace search term only one in a file, keep duplicaes (keep, default) or remove them (unique)
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

| Mode       | Description                                                                                                                                                    |
|:-----------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------|
| replace    | Replace search term inside one line with replacement.                                                                                                          |
| line       | Replace line (if matched term is inside) with replacement.                                                                                                     |
| lineinfile | Replace line (if matched term is inside) with replacement. If no match is found in the whole file the line will be appended to the bottom of the file.         |
| template   | Parse content as [golang template](https://golang.org/pkg/text/template/), arguments are available via `{{.Arg.Name}}` or environment vars via `{{.Env.Name}}` |


### Examples

| Command                                                            | Description                                                                                      |
|:-------------------------------------------------------------------|:-------------------------------------------------------------------------------------------------|
| `goreplace -s foobar -r barfoo file1 file2`                        | Replaces `foobar` to `barfoo` in file1 and file2                                                 |
| `goreplace --regex -s 'foo.*' -r barfoo file1 file2`               | Replaces the regex `foo.*` to `barfoo` in file1 and file2                                        |
| `goreplace --regex --ignore-case -s 'foo.*' -r barfoo file1 file2` | Replaces the regex `foo.*` (and ignore case) to `barfoo` in file1 and file2                      |
| `goreplace --mode=line -s 'foobar' -r barfoo file1 file2`          | Replaces all lines with content `foobar` to `barfoo` (whole line) in file1 and file2             |
| `goreplace -s 'foobar' -r barfoo --path=./ --path-pattern='*.txt'` | Replaces all lines with content `foobar` to `barfoo` (whole line) in *.txt files in current path |

### Example with golang templates

Configuration file `daemon.conf`:
```
<VirtualHost ...>
    ServerName {{.Env.SERVERNAME}}
    DocumentRoot {{.Env.DOCUMENTROOT}}
<VirtualHost>

```

Process file with:

```bash
SERVERNAME=www.foobar.example
DOCUMENTROOT=/var/www/foobar.example/
go-replace --mode=template daemon.conf
```

## Installation

```bash
GOREPLACE_VERSION=0.5.4 \
&& wget -O /usr/local/bin/go-replace https://github.com/webdevops/goreplace/releases/download/$GOREPLACE_VERSION/gr-64-linux \
&& chmod +x /usr/local/bin/go-replace
```

