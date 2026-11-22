package main

import (
	"errors"
	"fmt"
)

// Battery temperature limits in Celsius.
const (
	minTempC = -10
	maxTempC = 45
)

// ErrDuplicate means a frame's sequence number isn't newer than the last one stored.
var ErrDuplicate = errors.New("duplicate frame")

// LimitError is a reading outside its allowed range.
type LimitError struct {
	Field           string
	Value, Min, Max float64
}

// Error describes the reading and its range in a sentence.
func (e *LimitError) Error() string {
	return fmt.Sprintf("%s %.1f is outside %.1f to %.1f", e.Field, e.Value, e.Min, e.Max)
}

// TelemetryFrame is one decoded housekeeping frame.
type TelemetryFrame struct {
	SatID string
	Seq   uint32
	TempC float64
}

// validate returns every problem with f, given the last sequence number stored.
func validate(f TelemetryFrame, lastSeq uint32) error {
	if f.Seq <= lastSeq {
		return ErrDuplicate
	}
	if f.TempC < minTempC || f.TempC > maxTempC {
		return &LimitError{Field: "temp_c", Value: f.TempC, Min: minTempC, Max: maxTempC}
	}
	return nil
}

func main() {
	f := TelemetryFrame{SatID: "thistle", Seq: 4182, TempC: 900}
	fmt.Println(validate(f, 4182))
}
