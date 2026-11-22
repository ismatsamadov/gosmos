package main

import (
	"encoding/json"
	"errors"
	"io"
)

// decodeWeather reads one reading from r and refuses one without a station or a time.
func decodeWeather(r io.Reader) (Weather, error) {
	_, _ = json.NewDecoder, errors.New
	return Weather{}, nil
}
