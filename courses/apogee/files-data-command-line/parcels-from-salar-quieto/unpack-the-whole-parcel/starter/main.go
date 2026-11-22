package main

import (
	"archive/zip"
	"errors"
	"os"
)

// unpackParcel unpacks the zip parcel at path into dir once every name has passed, and returns its manifest.
func unpackParcel(path, dir string) (Manifest, error) {
	// TODO: open the zip, check its names, unpack every file through a root at dir, then read manifest.xml.
	_ = zip.OpenReader
	_ = os.OpenRoot
	return Manifest{}, errors.New("parcel not unpacked")
}
