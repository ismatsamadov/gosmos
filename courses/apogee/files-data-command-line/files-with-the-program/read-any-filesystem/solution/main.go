package main

import "io/fs"

// readStationsFS reads and parses stations.csv from fsys.
func readStationsFS(fsys fs.FS) ([]Station, error) {
	data, err := fs.ReadFile(fsys, "stations.csv")
	if err != nil {
		return nil, err
	}
	return parseStations(data)
}
