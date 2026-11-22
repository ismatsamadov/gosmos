package main

import (
	"archive/zip"
	"errors"
	"io"
)

// parcelManifest opens the zip parcel at path and returns its manifest, as readManifest checks it.
func parcelManifest(path string) (Manifest, error) {
	// TODO: open the zip, find manifest.xml among its files, read it whole and hand it to readManifest.
	_ = zip.OpenReader
	_ = io.ReadAll
	return Manifest{}, errors.New("parcel: no manifest.xml")
}
