package main

import (
	"archive/zip"
	"os"
)

// unpackFile writes the zipped file f below root, in the folder its name gives, refusing it past limit bytes.
func unpackFile(root *os.Root, f *zip.File, limit int64) error {
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	// TODO: make f's folder and create f through root, then copy src into it, at most limit bytes.
	return nil
}
