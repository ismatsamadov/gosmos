package main

import (
	"archive/zip"
	"path/filepath"
)

// checkNames returns an error for the first file whose name would lead out of the folder it is unpacked into.
func checkNames(files []*zip.File) error {
	// TODO: refuse the first file whose name filepath.IsLocal rejects, quoting the name.
	_ = filepath.IsLocal
	return nil
}
