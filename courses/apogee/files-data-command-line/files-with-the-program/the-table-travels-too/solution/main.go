package main

import _ "embed"

// stationsCSV is the station table, copied into gctl when it is built.
//
//go:embed stations.csv
var stationsCSV string

// stationTable returns the station table that gctl carries inside itself.
func stationTable() string {
	return stationsCSV
}
