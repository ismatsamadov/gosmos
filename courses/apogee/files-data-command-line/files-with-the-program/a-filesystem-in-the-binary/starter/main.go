package main

import (
	"embed"
	"path"
)

// bundled holds the templates folder, copied into gctl when it is built.
// TODO: embed the templates folder in this variable.
var bundled embed.FS

// loadTemplate returns the text of the built-in template templates/<name>.
func loadTemplate(name string) (string, error) {
	_ = path.Join
	return "", nil
}
