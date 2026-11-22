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

// LimitError is a reading outside its allowed range.
type LimitError struct {
	Field           string
	Value, Min, Max float64
}

// Error describes the reading and its range in a sentence.
func (e *LimitError) Error() string {
	return fmt.Sprintf("%s %.1f is outside %.1f to %.1f", e.Field, e.Value, e.Min, e.Max)
}

// checkTemp returns an error when the battery temperature c is outside the limits.
func checkTemp(c float64) error {
	if c < minTempC || c > maxTempC {
		return errors.New("temperature out of range")
	}
	return nil
}

// pagerDetail returns the field, value and range of the LimitError in err, or "" if there is none.
func pagerDetail(err error) string {
	return ""
}

func main() {
	err := fmt.Errorf("thistle seq 4190: %w", checkTemp(51))
	fmt.Println(err)
	fmt.Println(pagerDetail(err))
}
