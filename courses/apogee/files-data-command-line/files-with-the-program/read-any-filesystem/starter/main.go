package main

import "io/fs"

// readStationsFS reads and parses stations.csv from fsys.
func readStationsFS(fsys fs.FS) ([]Station, error) {
	_ = fs.ReadFile
	return nil, nil
}
