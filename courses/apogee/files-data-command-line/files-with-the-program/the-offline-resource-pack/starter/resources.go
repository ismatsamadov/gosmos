package main

import (
	"bytes"
	"embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strconv"
	"strings"
)

// resources holds the files gctl carries inside itself: the station table, the templates and the default calibration.
//
//go:embed stations.csv templates calibration.json
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

// Offsets is how far a dish points from where it should, in degrees.
type Offsets struct {
	AzOffsetDeg float64 `json:"az_offset_deg"`
	ElOffsetDeg float64 `json:"el_offset_deg"`
}

// Resources is what gctl needs at Salar Quieto: the station table and the van dish's offsets.
type Resources struct {
	Stations  []Station
	Reference Offsets
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

// readStationsFS reads and parses stations.csv from fsys.
func readStationsFS(fsys fs.FS) ([]Station, error) {
	data, err := fs.ReadFile(fsys, "stations.csv")
	if err != nil {
		return nil, err
	}
	return parseStations(data)
}

// parseOffsets reads calibration.json strictly: one object with known fields only.
func parseOffsets(data []byte) (Offsets, error) {
	var o Offsets
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&o); err != nil {
		return Offsets{}, fmt.Errorf("calibration.json: %w", err)
	}
	if err := dec.Decode(&json.RawMessage{}); !errors.Is(err, io.EOF) {
		return Offsets{}, errors.New("calibration.json: more than one value")
	}
	return o, nil
}
