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
	var env Envelope
	if err := decodeOne(r, &env); err != nil {
		return Report{}, err
	}
	switch env.Version {
	case 1:
		env.Kind = "image"
	case 2:
	default:
		return Report{}, fmt.Errorf("report %d: unsupported version %d", env.ID, env.Version)
	}
	payload, err := decodePayload(env.Kind, env.Payload)
	if err != nil {
		return Report{}, fmt.Errorf("report %d: %w", env.ID, err)
	}
	return Report{ID: env.ID, Kind: env.Kind, Payload: payload}, nil
}
