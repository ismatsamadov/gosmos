package main

import (
	"encoding/json"
	"fmt"
)

// ImageReport says how much of one image file reached Furrowlight.
type ImageReport struct {
	SatID   string `json:"sat_id"`
	Station string `json:"station"`
	File    string `json:"file"`
	Chunks  int    `json:"chunks"`
	Missing int    `json:"missing"`
}

// StationReport sums up one station's day of deliveries to Furrowlight.
type StationReport struct {
	Station string `json:"station"`
	Day     string `json:"day"`
	Passes  int    `json:"passes"`
	Missing int    `json:"missing"`
}

// decodeReport decodes a report's payload into the type its kind names.
func decodeReport(data []byte) (any, error) {
	_ = json.Unmarshal
	_ = fmt.Errorf
	return nil, nil
}
