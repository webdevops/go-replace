package main

import (
	"bufio"
	"bytes"
	"regexp"
)

// check if string is contained in an array
func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

// Checks if there is a match in content, based on search options
func searchMatch(content string, changeset changeset) bool {
	return changeset.Search.MatchString(content)
}

// Replace text in whole content based on search options
func replaceText(content string, changeset changeset) string {
	// --regex-backrefs
	if opts.RegexBackref {
		return changeset.Search.ReplaceAllString(content, changeset.Replace)
	} else {
		return changeset.Search.ReplaceAllLiteralString(content, changeset.Replace)
	}
}

func handleLineInFile(changesets []changeset, buffer bytes.Buffer) (*bytes.Buffer, bool) {
	var (
		line              string
		writeBufferToFile bool
	)

	for _, changeset := range changesets {
		if !changeset.MatchFound {
			// just add line to file
			line = changeset.Replace + "\n"

			// remove backrefs (no match)
			if opts.RegexBackref {
				line = regexp.MustCompile(`\$[0-9]+`).ReplaceAllLiteralString(line, "")
			}

			// --lineinfile-before
			// --lineinfile-after
			if opts.LineinfileBefore != "" || opts.LineinfileAfter != "" {
				var matchFinder *regexp.Regexp

				if opts.LineinfileBefore != "" {
					matchFinder = regexp.MustCompile(opts.LineinfileBefore)
				} else {
					matchFinder = regexp.MustCompile(opts.LineinfileAfter)
				}

				var bufferCopy bytes.Buffer

				scanner := bufio.NewScanner(&buffer)
				for scanner.Scan() {
					originalLine := scanner.Text()

					if matchFinder.MatchString(originalLine) {
						writeBufferToFile = true

						if opts.LineinfileBefore != "" {
							bufferCopy.WriteString(line)
						}

						bufferCopy.WriteString(originalLine + "\n")

						if opts.LineinfileAfter != "" {
							bufferCopy.WriteString(line)
						}
					} else {
						bufferCopy.WriteString(originalLine + "\n")
					}
				}

				buffer.Reset()
				buffer.WriteString(bufferCopy.String())
			} else {
				buffer.WriteString(line)
				writeBufferToFile = true
			}
		}
	}

	return &buffer, writeBufferToFile
}
