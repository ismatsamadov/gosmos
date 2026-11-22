package main

import (
	"errors"
	"io/fs"
)

// loadResources reads the station table from fsys, and calibration.json from override or else from fsys.
func loadResources(fsys, override fs.FS) (Resources, error) {
	stations, err := readStationsFS(fsys)
	if err != nil {
		return Resources{}, err
	}
	data, err := fs.ReadFile(override, "calibration.json")
	if errors.Is(err, fs.ErrNotExist) {
		data, err = fs.ReadFile(fsys, "calibration.json")
	}
	if err != nil {
		return Resources{}, err
	}
	reference, err := parseOffsets(data)
	if err != nil {
		return Resources{}, err
	}
	return Resources{Stations: stations, Reference: reference}, nil
}
