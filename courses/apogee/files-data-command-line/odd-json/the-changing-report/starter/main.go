package main

import (
	"fmt"
	"io"
)

// Report is a delivery report in today's form, whatever version it came in.
type Report struct {
	ID      int64
	Kind    string
	Payload any // *ImageReport or *StationReport
}

// normalizeReport reads one report of a known version strictly and returns it in today's form.
func normalizeReport(r io.Reader) (Report, error) {
	_ = fmt.Errorf
	return Report{}, nil
}
