package main

import (
	"archive/zip"
	"fmt"
	"path/filepath"
)

// checkNames returns an error for the first file whose name would lead out of the folder it is unpacked into.
func checkNames(files []*zip.File) error {
	for _, f := range files {
		if !filepath.IsLocal(f.Name) {
			return fmt.Errorf("unsafe name %q", f.Name)
		}
	}
	return nil
}
