package main

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	sprig "github.com/Masterminds/sprig"
)

type templateData struct {
	Arg map[string]string
	Env map[string]string
}

func createTemplate() *template.Template {
	tmpl := template.New("base")
	tmpl.Funcs(sprig.TxtFuncMap())
	tmpl.Option("missingkey=zero")

	return tmpl
}

func parseContentAsTemplate(templateContent string, changesets []changeset) bytes.Buffer {
	var content bytes.Buffer
	data := generateTemplateData(changesets)
	tmpl, err := createTemplate().Parse(templateContent)
	if err != nil {
		logFatalErrorAndExit(err, 1)
	}

	err = tmpl.Execute(&content, &data)
	if err != nil {
		logFatalErrorAndExit(err, 1)
	}

	return content
}

func generateTemplateData(changesets []changeset) templateData {
	// init
	var ret templateData
	ret.Arg = make(map[string]string)
	ret.Env = make(map[string]string)

	// add changesets
	for _, changeset := range changesets {
		ret.Arg[changeset.SearchPlain] = changeset.Replace
	}

	// add env variables
	for _, e := range os.Environ() {
		split := strings.SplitN(e, "=", 2)
		envKey, envValue := split[0], split[1]
		ret.Env[envKey] = envValue
	}

	return ret
}
