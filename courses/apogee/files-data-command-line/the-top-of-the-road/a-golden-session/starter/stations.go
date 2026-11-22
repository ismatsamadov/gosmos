package main

import (
	_ "embed"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
)

// stationsCSV is the station table, copied into gctl when it is built.
//
//go:embed stations.csv
var stationsCSV string

// Band is a radio band.
type Band string

// Station is one of Apogee's ground stations.
type Station struct {
	Code    string // "YST": three capital letters
	Name    string // "Ystmark"
	MaskDeg float64
	Bands   []Band
}

// loadStations reads the built-in station table, by station code.
func loadStations() (map[string]Station, error) {
	rows, err := csv.NewReader(strings.NewReader(stationsCSV)).ReadAll()
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
