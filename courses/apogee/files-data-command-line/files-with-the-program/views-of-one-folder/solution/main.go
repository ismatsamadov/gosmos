package main

import "io/fs"

// templateFS returns a view of fsys whose root is its templates folder.
func templateFS(fsys fs.FS) (fs.FS, error) {
	if _, err := fs.Stat(fsys, "templates"); err != nil {
		return nil, err
	}
	return fs.Sub(fsys, "templates")
}
