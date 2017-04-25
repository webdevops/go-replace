package main

import (
    "fmt"
    "bytes"
    "io/ioutil"
    "path/filepath"
    "bufio"
    "os"
    "regexp"
)

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

// Write content to file
func writeContentToFile(fileitem fileitem, content bytes.Buffer) (string, bool) {
    // --dry-run
    if opts.DryRun {
        return content.String(), true
    } else {
        var err error
        err = ioutil.WriteFile(fileitem.Output, content.Bytes(), 0644)
        if err != nil {
            panic(err)
        }

        return fmt.Sprintf("%s found and replaced match\n", fileitem.Path), true
    }
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
