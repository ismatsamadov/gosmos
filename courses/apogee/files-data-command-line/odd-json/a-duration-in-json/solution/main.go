package main

import (
	"encoding/json"
	"time"
)

// Delay is a timeout in the settings file, written as duration text such as "30s".
type Delay time.Duration

// MarshalJSON writes d as duration text, such as "30s".
func (d Delay) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// UnmarshalJSON reads duration text, such as "2m", into d.
func (d *Delay) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}
	v, err := time.ParseDuration(text)
	if err != nil {
		return err
	}
	*d = Delay(v)
	return nil
}

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
