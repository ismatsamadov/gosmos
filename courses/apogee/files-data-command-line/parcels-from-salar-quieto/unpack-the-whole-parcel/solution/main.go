package main

import (
	"archive/zip"
	"os"
)

// unpackParcel unpacks the zip parcel at path into dir once every name has passed, and returns its manifest.
func unpackParcel(path, dir string) (Manifest, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return Manifest{}, err
	}
	defer zr.Close()
	if err := checkNames(zr.File); err != nil {
		return Manifest{}, err
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return Manifest{}, err
	}
	defer root.Close()
	for _, f := range zr.File {
		if err := unpackFile(root, f, maxFileBytes); err != nil {
			return Manifest{}, err
		}
	}
	data, err := root.ReadFile("manifest.xml")
	if err != nil {
		return Manifest{}, err
	}
	return readManifest(data)
}
