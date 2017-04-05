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


Testing ignoring missing arguments:
  $ goreplace -s foobar -r ___xxx --ignore-empty



Testing replace mode:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the last line
  > EOF
  $ goreplace -s foobar -r ___xxx test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  this is the third ___xxx line
  this is the last line

Testing replace mode with multiple matches:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the foobar forth foobar line
  > this is the last line
  > EOF
  $ goreplace -s foobar -r ___xxx test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  this is the third ___xxx line
  this is the ___xxx forth ___xxx line
  this is the last line

Testing replace mode without match:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the foobar forth foobar line
  > this is the last line
  > EOF
  $ goreplace -s barfoo -r ___xxx test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  this is the third foobar line
  this is the foobar forth foobar line
  this is the last line

Testing replace mode with regex:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the last line
  > EOF
  $ goreplace --regex -s 'f[o]+b[a]*r' -r ___xxx test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  this is the third ___xxx line
  this is the last line

Testing replace mode with regex:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the last line
  > EOF
  $ goreplace --regex --regex-backrefs -s 'f[o]+(b[a]*r)' -r '___$1' test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  this is the third ___bar line
  this is the last line

Testing replace mode with regex and case-insensitive:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the last line
  > EOF
  $ goreplace --regex --regex-backrefs -s 'F[O]+(b[a]*r)' -r '___$1' --case-insensitive test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  this is the third ___bar line
  this is the last line


Testing line mode:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the last line
  > EOF
  $ goreplace --mode=line -s foobar -r ___xxx test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  ___xxx
  this is the last line

Testing line mode with multiple matches:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the foobar forth foobar line
  > this is the last line
  > EOF
  $ goreplace --mode=line -s foobar -r ___xxx test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  ___xxx
  ___xxx
  this is the last line

Testing line mode with multiple matches and --once:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the foobar forth foobar line
  > this is the last line
  > EOF
  $ goreplace --mode=line -s foobar -r ___xxx --once test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  ___xxx
  this is the foobar forth foobar line
  this is the last line

Testing line mode with multiple matches and --once-remove-match:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the foobar forth foobar line
  > this is the last line
  > EOF
  $ goreplace --mode=line -s foobar -r ___xxx --once-remove-match test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  ___xxx
  this is the last line




Testing lineinfile mode:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the last line
  > EOF
  $ goreplace --mode=lineinfile -s foobar -r ___xxx test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  ___xxx
  this is the last line

Testing lineinfile mode with multiple matches:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the foobar forth foobar line
  > this is the last line
  > EOF
  $ goreplace --mode=lineinfile -s foobar -r ___xxx test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  ___xxx
  ___xxx
  this is the last line

Testing lineinfile mode with multiple matches and --once:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the foobar forth foobar line
  > this is the last line
  > EOF
  $ goreplace --mode=lineinfile -s foobar -r ___xxx --once test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  ___xxx
  this is the foobar forth foobar line
  this is the last line

Testing lineinfile mode with multiple matches and --once-remove-match:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the foobar forth foobar line
  > this is the last line
  > EOF
  $ goreplace --mode=lineinfile -s foobar -r ___xxx --once-remove-match test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  ___xxx
  this is the last line

Testing lineinfile mode without match:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the foobar forth foobar line
  > this is the last line
  > EOF
  $ goreplace --mode=lineinfile -s barfoo -r ___xxx --once-remove-match test.txt
  $ cat test.txt
  this is a testline
  this is the second line
  this is the third foobar line
  this is the foobar forth foobar line
  this is the last line
  ___xxx



Testing replace mode with path option:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the foobar forth foobar line
  > this is the last line
  > EOF
  $ mkdir -p testing/sub1/subsub testing/sub2/subsub
  $ cp test.txt testing/sub1/subsub/test1.txt
  $ cp test.txt testing/sub1/subsub/test2.txt
  $ cp test.txt testing/sub1/test3.txt
  $ cp test.txt testing/sub2/subsub/test4.txt
  $ cp test.txt testing/sub2/subsub/test5.txt
  $ cp test.txt testing/sub2/original.md
  $ goreplace -s foobar -r barfoo --path=./testing --path-pattern='*.txt'
  $ cat testing/sub1/subsub/test1.txt
  this is a testline
  this is the second line
  this is the third barfoo line
  this is the barfoo forth barfoo line
  this is the last line
  $ cat testing/sub1/subsub/test2.txt
  this is a testline
  this is the second line
  this is the third barfoo line
  this is the barfoo forth barfoo line
  this is the last line
  $ cat testing/sub1/test3.txt
  this is a testline
  this is the second line
  this is the third barfoo line
  this is the barfoo forth barfoo line
  this is the last line
  $ cat testing/sub2/subsub/test4.txt
  this is a testline
  this is the second line
  this is the third barfoo line
  this is the barfoo forth barfoo line
  this is the last line
  $ cat testing/sub2/subsub/test5.txt
  this is a testline
  this is the second line
  this is the third barfoo line
  this is the barfoo forth barfoo line
  this is the last line
  $ cat testing/sub2/original.md
  this is a testline
  this is the second line
  this is the third foobar line
  this is the foobar forth foobar line
  this is the last line
