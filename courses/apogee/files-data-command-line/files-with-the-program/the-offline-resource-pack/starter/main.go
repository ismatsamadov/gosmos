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
	// TODO: read calibration.json from override, or from fsys when override has none, and parse it.
	_ = errors.Is
	return Resources{Stations: stations}, nil
}
