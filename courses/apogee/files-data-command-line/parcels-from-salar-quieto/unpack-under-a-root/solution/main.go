package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

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
