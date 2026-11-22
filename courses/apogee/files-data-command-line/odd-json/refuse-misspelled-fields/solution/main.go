package main

import (
	"encoding/json"
	"io"
)

// Settings is gctl's settings file.
type Settings struct {
	Station    string  `json:"station"`
	Day        string  `json:"day"`
	MinElevDeg float64 `json:"min_elev_deg"`
}

// strictConfig decodes one Settings value from r, refusing fields Settings doesn't have.
func strictConfig(r io.Reader) (Settings, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var s Settings
	if err := dec.Decode(&s); err != nil {
		return Settings{}, err
	}
	return s, nil
}
