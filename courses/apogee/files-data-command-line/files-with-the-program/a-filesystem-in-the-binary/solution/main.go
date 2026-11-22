package main

import (
	"embed"
	"path"
)

// bundled holds the templates folder, copied into gctl when it is built.
//
//go:embed templates
var bundled embed.FS

// loadTemplate returns the text of the built-in template templates/<name>.
func loadTemplate(name string) (string, error) {
	data, err := bundled.ReadFile(path.Join("templates", name))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
