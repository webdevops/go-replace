Go Replace tests:

  $ go build -o goreplace "$TESTDIR/../goreplace.go"
  $ CURRENT=$(pwd)
  $ alias goreplace="$CURRENT/goreplace"

Usage:

  $ goreplace -h
  Usage:
    goreplace
  
  Application Options:
    -m, --mode=[replace|line|lineinfile] replacement mode - replace: replace
                                         match with term; line: replace line with
                                         term; lineinfile: replace line with term
                                         or if not found append to term to file
                                         (default: replace)
    -s, --search=                        search term
    -r, --replace=                       replacement term
    -i, --case-insensitive               ignore case of pattern to match upper
                                         and lowercase characters
        --once                           replace search term only one in a file
        --once-remove-match              replace search term only one in a file
                                         and also don't keep matching lines (for
                                         line and lineinfile mode)
        --regex                          treat pattern as regex
        --regex-backrefs                 enable backreferences in replace term
        --regex-posix                    parse regex term as POSIX regex
        --path=                          use files in this path
        --path-pattern=                  file pattern (* for wildcard, only
                                         basename of file)
        --path-regex=                    file pattern (regex, full path)
        --ignore-empty                   ignore empty file list, otherwise this
                                         will result in an error
    -v, --verbose                        verbose mode
        --dry-run                        dry run mode
    -V, --version                        show version and exit
    -h, --help                           show this help message
  [1]
  $ goreplace -V
  goreplace version [0-9]+.[0-9]+.[0-9]+ (re)
