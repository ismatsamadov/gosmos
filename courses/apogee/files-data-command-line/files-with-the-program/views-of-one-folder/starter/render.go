package main

import (
	"embed"
	"io/fs"
	"strings"
)

// resources holds the files gctl carries inside itself: the station table and the templates.
//
//go:embed stations.csv templates
var resources embed.FS

// render reads the template name from templates and fills in its {field} placeholders.
func render(templates fs.FS, name string, fields map[string]string) (string, error) {
	data, err := fs.ReadFile(templates, name)
	if err != nil {
		return "", err
	}
	text := string(data)
	for field, value := range fields {
		text = strings.ReplaceAll(text, "{"+field+"}", value)
	}
	return text, nil
}
