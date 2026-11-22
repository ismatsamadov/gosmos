package main

import (
	"bytes"
	"embed"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// resources holds the files gctl carries inside itself.
//
//go:embed stations.csv
var resources embed.FS

// Band is a radio band.
type Band string

// Station is one of Apogee's ground stations.
type Station struct {
	Code    string // "YST": three capital letters
	Name    string // "Ystmark"
	MaskDeg float64
	Bands   []Band
}

// parseStations reads a station table: its header line, then one station per line.
func parseStations(data []byte) ([]Station, error) {
	rows, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("station table: %w", err)
	}
	if len(rows) == 0 || strings.Join(rows[0], ",") != "code,name,lat_deg,lon_deg,mask_deg,bands" {
		return nil, errors.New("station table: the first line is not the header")
	}
	var stations []Station
	for _, row := range rows[1:] {
		mask, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return nil, fmt.Errorf("station %s: mask: %w", row[0], err)
		}
		s := Station{Code: row[0], Name: row[1], MaskDeg: mask}
		for _, band := range strings.Fields(row[5]) {
			s.Bands = append(s.Bands, Band(band))
		}
		stations = append(stations, s)
	}
	return stations, nil
}
