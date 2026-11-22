package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// loadConfig reads one settings object strictly over the defaults and checks its day.
func loadConfig(r io.Reader) (Settings, error) {
	settings := Settings{Station: "SQT", MinElevDeg: 3, Timeout: Delay(30 * time.Second)}
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&settings); err != nil {
		return Settings{}, err
	}
	if err := dec.Decode(&json.RawMessage{}); !errors.Is(err, io.EOF) {
		return Settings{}, errors.New("more than one JSON value")
	}
	if _, err := time.Parse(time.DateOnly, settings.Day); err != nil {
		return Settings{}, fmt.Errorf("day %q: want a date such as 2027-10-21", settings.Day)
	}
	return settings, nil
}
