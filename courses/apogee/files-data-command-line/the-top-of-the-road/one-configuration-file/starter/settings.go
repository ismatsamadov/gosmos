package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Delay is a duration that settings.json writes as text, such as "30s" or "2m".
type Delay time.Duration

// UnmarshalJSON reads duration text, refusing numbers and durations that are not positive.
func (d *Delay) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return fmt.Errorf("timeout: want text such as \"30s\": %w", err)
	}
	v, err := time.ParseDuration(text)
	if err != nil {
		return fmt.Errorf("timeout: %w", err)
	}
	if v <= 0 {
		return fmt.Errorf("timeout %q: want more than zero", text)
	}
	*d = Delay(v)
	return nil
}

// String writes the delay as duration text, such as 30s.
func (d Delay) String() string {
	return time.Duration(d).String()
}

// Settings is gctl's settings file, with defaults for what it leaves out.
type Settings struct {
	Station    string  `json:"station"`
	Day        string  `json:"day"`
	MinElevDeg float64 `json:"min_elev_deg"`
	Timeout    Delay   `json:"timeout"`
}
