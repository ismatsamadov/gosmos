package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Envelope is one Furrowlight delivery report, its payload not yet decoded.
type Envelope struct {
	Version int             `json:"version"`
	Kind    string          `json:"kind"`
	ID      int64           `json:"id"`
	Payload json.RawMessage `json:"payload"`
}

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

// decodeOne decodes exactly one JSON value from r into v, refusing unknown fields.
func decodeOne(r io.Reader, v any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return err
	}
	if _, err := dec.Token(); err != io.EOF {
		return errors.New("unexpected data after the JSON value")
	}
	return nil
}

// decodePayload decodes a payload strictly into the report type its kind names.
func decodePayload(kind string, raw json.RawMessage) (any, error) {
	var report any
	switch kind {
	case "image":
		report = &ImageReport{}
	case "station":
		report = &StationReport{}
	default:
		return nil, fmt.Errorf("unknown report kind %q", kind)
	}
	if len(raw) == 0 || string(raw) == "null" {
		return nil, errors.New("no payload")
	}
	if err := decodeOne(bytes.NewReader(raw), report); err != nil {
		return nil, err
	}
	return report, nil
}
