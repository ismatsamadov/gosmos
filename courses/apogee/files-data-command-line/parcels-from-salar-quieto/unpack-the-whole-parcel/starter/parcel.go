package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// maxFileBytes is the most gctl unpacks of any one file in a parcel: 4 MiB.
const maxFileBytes = 4 << 20

// checkNames returns an error for the first file whose name would lead out of the folder it is unpacked into.
func checkNames(files []*zip.File) error {
	for _, f := range files {
		if !filepath.IsLocal(f.Name) {
			return fmt.Errorf("unsafe name %q", f.Name)
		}
	}
	return nil
}

// unpackFile writes the zipped file f below root, in the folder its name gives, refusing it past limit bytes.
func unpackFile(root *os.Root, f *zip.File, limit int64) error {
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	if err := root.MkdirAll(filepath.Dir(f.Name), 0o700); err != nil {
		return err
	}
	dst, err := root.Create(f.Name)
	if err != nil {
		return err
	}
	defer dst.Close()
	n, err := io.Copy(dst, io.LimitReader(src, limit+1))
	if err != nil {
		return err
	}
	if n > limit {
		return fmt.Errorf("%s: more than %d bytes", f.Name, limit)
	}
	return dst.Close()
}
