package main

import (
	"encoding/json"
	"errors"
	"io"
)

// ErrTrailing means a settings file holds something after its one object.
var ErrTrailing = errors.New("settings: unexpected data after the object")

// Settings is gctl's settings file.
type Settings struct {
	Station    string  `json:"station"`
	Day        string  `json:"day"`
	MinElevDeg float64 `json:"min_elev_deg"`
}

// decodeOne decodes exactly one Settings object from r, refusing unknown fields and anything after it.
func decodeOne(r io.Reader) (Settings, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var s Settings
	if err := dec.Decode(&s); err != nil {
		return Settings{}, err
	}
	if _, err := dec.Token(); err != io.EOF {
		return Settings{}, ErrTrailing
	}
	return s, nil
}
