package main

import (
	"encoding/json"
	"time"
)

// Delay is a timeout in the settings file, written as duration text such as "30s".
type Delay time.Duration

// TODO: give Delay a MarshalJSON method and *Delay an UnmarshalJSON method that use duration text.

// Settings is gctl's settings file.
type Settings struct {
	Station    string  `json:"station"`
	Day        string  `json:"day"`
	MinElevDeg float64 `json:"min_elev_deg"`
	Timeout    Delay   `json:"timeout"`
}

// saveSettings encodes s as the settings file's JSON.
func saveSettings(s Settings) ([]byte, error) { return json.Marshal(s) }

// loadSettings decodes the settings file's JSON.
func loadSettings(data []byte) (Settings, error) {
	var s Settings
	err := json.Unmarshal(data, &s)
	return s, err
}
