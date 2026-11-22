package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strconv"
	"strings"
)

// Band is a radio band.
type Band string

// Station is one of Apogee's ground stations, as the station table lists it.
type Station struct {
	Code    string
	Name    string
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

// stowNotice fills the stow notice template for a station of the table, given by its code.
func stowNotice(fsys fs.FS, code, reason string) (string, error) {
	table, err := fs.ReadFile(fsys, "stations.csv")
	if err != nil {
		return "", err
	}
	stations, err := parseStations(table)
	if err != nil {
		return "", err
	}
	if !slices.ContainsFunc(stations, func(s Station) bool { return s.Code == code }) {
		return "", fmt.Errorf("unknown station %q", code)
	}
	notice, err := fs.ReadFile(fsys, "templates/notices/stow.txt")
	if err != nil {
		return "", err
	}
	return strings.NewReplacer("{station}", code, "{reason}", reason).Replace(string(notice)), nil
}
