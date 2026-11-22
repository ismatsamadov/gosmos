package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

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
