Go Replace tests:

  $ go build -o goreplace "$TESTDIR/../goreplace.go"
  $ CURRENT=$(pwd)
  $ alias goreplace="$CURRENT/goreplace"

Usage:

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

Testing replace mode with multiple changesets:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second barfoo line
  > this is the third foobar line
  > this is the last oofrab line
  > EOF
  $ goreplace -s foobar -r 111 -s barfoo -r 222 -s oofrab -r 333 test.txt
  $ cat test.txt
  this is a testline
  this is the second 222 line
  this is the third 111 line
  this is the last 333 line

Testing replace mode with stdin:

  $ cat > test.txt <<EOF
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the last line
  > EOF
  $ cat test.txt | goreplace -s foobar -r ___xxx --stdin
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


Testing template mode:

  $ cat > test.txt <<EOF
  > {{23 -}} < {{- 45}}
  > {{.Arg.Foobar}}
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the last line
  > EOF
  $ goreplace --mode=template -s Foobar -r ___xxx test.txt
  $ cat test.txt
  23<45
  ___xxx
  this is a testline
  this is the second line
  this is the third foobar line
  this is the last line

Testing template mode with only env:

  $ cat > test.txt <<EOF
  > {{23 -}} < {{- 45}}
  > {{.Env.Foobar}}
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the last line
  > EOF
  $ Foobar=barfoo goreplace --mode=template test.txt
  $ cat test.txt
  23<45
  barfoo
  this is a testline
  this is the second line
  this is the third foobar line
  this is the last line

Testing template mode with only env and empty var:

  $ cat > test.txt <<EOF
  > {{23 -}} < {{- 45}}
  > begin{{.Env.Foobar}}end
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the last line
  > EOF
  $ Foobar= goreplace --mode=template test.txt
  $ cat test.txt
  23<45
  beginend
  this is a testline
  this is the second line
  this is the third foobar line
  this is the last line


Testing template mode:

  $ cat > test.txt <<EOF
  > {{23 -}} < {{- 45}}
  > {{.Arg.Foobar}}
  > this is a testline
  > this is the second line
  > this is the third foobar line
  > this is the last line
  > EOF
  $ cat test.txt | goreplace --mode=template --stdin -s Foobar -r ___xxx test.txt
  23<45
  ___xxx
  this is a testline
  this is the second line
  this is the third foobar line
  this is the last line
