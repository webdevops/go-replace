// The implementation is based of
// https://github.com/gliderlabs/sigil/blob/master/builtin/builtin.go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"io"
)

type NamedReader struct {
	io.Reader
	Name string
}

// Merge template.FuncMaps into a single map
// First entry win's
func Merge(ms ...template.FuncMap) template.FuncMap {
	res := make(template.FuncMap)
	for _, m := range ms {
	srcMap:
		for k, v := range m {
			// Check if (k,v) was added before:
			_, ok := res[k]
			if ok {
				continue srcMap
			}
			res[k] = v
		}
	}
	return res
}

// Extract data from the interface
func String(in interface{}) (string, string, bool) {
	switch obj := in.(type) {
	case string:
		return obj, "", true
	case NamedReader:
		data, err := ioutil.ReadAll(obj)
		if err != nil {
			panic(err)
		}
		return string(data), obj.Name, true
	case fmt.Stringer:
		return obj.String(), "", true
	default:
		return "", "", false
	}
}

func LookPath(file string) (string, error) {
	if strings.HasPrefix(file, "/") {
		return file, nil
	}
	cwd, _ := os.Getwd()
	search := []string{cwd}
	for _, path := range search {
		filepath := filepath.Join(path, file)
		if _, err := os.Stat(filepath); err == nil {
			return filepath, nil
		}
	}
	return "", fmt.Errorf("Not found in path: %s", file)
}

func builtInFunctions() template.FuncMap {
	return template.FuncMap{
		"file":   File,
	}
}

func read(file interface{}) ([]byte, error) {
	reader, ok := file.(NamedReader)
	if ok {
		return ioutil.ReadAll(reader)
	}
	path, _, ok := String(file)
	if !ok {
		return []byte{}, fmt.Errorf("file must be stream or path string")
	}
	fp, err := LookPath(path)
	if err != nil {
		return []byte{}, err
	}
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func File(filename interface{}) (interface{}, error) {
	str, err := read(filename)
	return string(str), err
}
