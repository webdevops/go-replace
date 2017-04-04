# goreplace

[![GitHub release](https://img.shields.io/github/release/webdevops/goreplace.svg)](https://github.com/webdevops/goreplace/releases)
[![license](https://img.shields.io/github/license/webdevops/goreplace.svg)](https://github.com/webdevops/goreplace/blob/master/LICENSE)

Cli utility for replacing text in files, written in golang and compiled for usage in Docker images

Inspired by https://github.com/piranha/goreplace

## Usage

```
Usage:
  goreplace

Application Options:
  -s, --search=         search term
  -r, --replace=        replacement term
  -i, --ignore-case     ignore pattern case
      --replace-line    replace whole line instead of only match
      --regex           treat pattern as regex
      --regex-backrefs  enable backreferences in replace term
      --regex-posix     parse regex term as POSIX regex
      --path=           use files in this path
      --path-pattern=   file pattern (* for wildcard, only basename of file)
      --path-regex=     file pattern (regex, full path)
      --ignore-empty    ignore empty file list, otherwise this will result in an error
  -v, --verbose         verbose mode
      --dry-run         dry run mode
  -V, --version         show version and exit
  -h, --help            show this help message
```

### Examples

| Command                                                            | Description                                                                                      |
|:-------------------------------------------------------------------|:-------------------------------------------------------------------------------------------------|
| `goreplace -s foobar -r barfoo file1 file2`                        | Replaces `foobar` to `barfoo` in file1 and file2                                                 |
| `goreplace --regex -s 'foo.*' -r barfoo file1 file2`               | Replaces the regex `foo.*` to `barfoo` in file1 and file2                                        |
| `goreplace --regex --ignore-case -s 'foo.*' -r barfoo file1 file2` | Replaces the regex `foo.*` (and ignore case) to `barfoo` in file1 and file2                      |
| `goreplace --replace-line -s 'foobar' -r barfoo file1 file2`       | Replaces all lines with content `foobar` to `barfoo` (whole line) in file1 and file2             |
| `goreplace -s 'foobar' -r barfoo --path=./ --path-pattern='*.txt'` | Replaces all lines with content `foobar` to `barfoo` (whole line) in *.txt files in current path |


## Installation

```bash
GOREPLACE_VERSION=0.3.1 \
&& wget -O /usr/local/bin/go-replace https://github.com/webdevops/goreplace/releases/download/$GOREPLACE_VERSION/gr-64-linux \
&& chmod +x /usr/local/bin/go-replace
```

