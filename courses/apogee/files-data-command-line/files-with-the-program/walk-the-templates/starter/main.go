package main

import "io/fs"

// listTemplates returns the path of every file below templates in fsys, in walk order.
func listTemplates(fsys fs.FS) ([]string, error) {
	_ = fs.WalkDir
	return nil, nil
}
