package main

import "io/fs"

// templateFS returns a view of fsys whose root is its templates folder.
func templateFS(fsys fs.FS) (fs.FS, error) {
	_ = fs.Sub
	return nil, nil
}
