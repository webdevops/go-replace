package main

import (
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
func searchMatch(content string, changeset changeset) (bool) {
    if changeset.Search.MatchString(content) {
        return true
    }

    return false
}

// Replace text in whole content based on search options
func replaceText(content string, changeset changeset) (string) {
    // --regex-backrefs
    if opts.RegexBackref {
        return changeset.Search.ReplaceAllString(content, changeset.Replace)
    } else {
        return changeset.Search.ReplaceAllLiteralString(content, changeset.Replace)
    }
}
