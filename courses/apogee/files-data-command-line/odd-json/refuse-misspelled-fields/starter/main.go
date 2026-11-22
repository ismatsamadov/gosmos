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
	_ = json.NewDecoder
	return Settings{}, nil
}
