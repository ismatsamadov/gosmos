package main

import (
	"encoding/json"
	"errors"
	"io"
)

// decodeWeather reads one reading from r and refuses one without a station or a time.
func decodeWeather(r io.Reader) (Weather, error) {
	var w Weather
	if err := json.NewDecoder(r).Decode(&w); err != nil {
		return Weather{}, err
	}
	if w.Station == "" {
		return Weather{}, errors.New(`reading has no "station"`)
	}
	if w.At.IsZero() {
		return Weather{}, errors.New(`reading has no "at"`)
	}
	return w, nil
}
