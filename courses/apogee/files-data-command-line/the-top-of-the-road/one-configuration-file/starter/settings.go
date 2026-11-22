package main

import (
	"encoding/json"
	"time"
)

// Delay is a timeout in the settings file, written as duration text such as "30s".
type Delay time.Duration

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

// String writes the delay as duration text, such as 30s.
func (d Delay) String() string {
	return time.Duration(d).String()
}

// Settings is gctl's settings file.
type Settings struct {
	Station    string  `json:"station"`
	Day        string  `json:"day"`
	MinElevDeg float64 `json:"min_elev_deg"`
	Timeout    Delay   `json:"timeout"`
}
