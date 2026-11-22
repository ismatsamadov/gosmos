package main

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
)

// stationsCSV is the station table gctl carries inside itself.
//
//go:embed stations.csv
var stationsCSV []byte

// Station is one of Apogee's ground stations, as the station table lists it.
type Station struct {
	Code string
	Name string
}

// loadStations reads the built-in station table, by station code.
func loadStations() (map[string]Station, error) {
	rows, err := csv.NewReader(bytes.NewReader(stationsCSV)).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("station table: %w", err)
	}
	if len(rows) == 0 || strings.Join(rows[0], ",") != "code,name,lat_deg,lon_deg,mask_deg,bands" {
		return nil, errors.New("station table: the first line is not the header")
	}
	stations := map[string]Station{}
	for _, row := range rows[1:] {
		stations[row[0]] = Station{Code: row[0], Name: row[1]}
	}
	return stations, nil
}
