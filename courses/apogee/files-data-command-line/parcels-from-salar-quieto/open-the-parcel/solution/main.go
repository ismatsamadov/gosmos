package main

import (
	"archive/zip"
	"errors"
	"io"
)

// parcelManifest opens the zip parcel at path and returns its manifest, as readManifest checks it.
func parcelManifest(path string) (Manifest, error) {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return Manifest{}, err
	}
	defer zr.Close()
	for _, f := range zr.File {
		if f.Name != "manifest.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return Manifest{}, err
		}
		defer rc.Close()
		data, err := io.ReadAll(rc)
		if err != nil {
			return Manifest{}, err
		}
		return readManifest(data)
	}
	return Manifest{}, errors.New("parcel: no manifest.xml")
}
