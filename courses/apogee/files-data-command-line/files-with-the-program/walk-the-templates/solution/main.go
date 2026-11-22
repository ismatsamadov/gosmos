package main

import "io/fs"

// listTemplates returns the path of every file below templates in fsys, in walk order.
func listTemplates(fsys fs.FS) ([]string, error) {
	var paths []string
	err := fs.WalkDir(fsys, "templates", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			paths = append(paths, p)
		}
		return nil
	})
	return paths, err
}
