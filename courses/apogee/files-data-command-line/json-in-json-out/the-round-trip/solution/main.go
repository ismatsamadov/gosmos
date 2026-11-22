package main

import (
	"encoding/json"
	"errors"
)

// confirm checks one booking from the export and returns it as JSON for the reply, its ID unchanged.
func confirm(data []byte) ([]byte, error) {
	var b Booking
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}
	if b.ID == "" {
		return nil, errors.New("booking has no id")
	}
	return json.Marshal(b)
}
