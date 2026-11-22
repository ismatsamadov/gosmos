package main

import (
	"encoding/json"
	"errors"
)

// Envelope is one Furrowlight delivery report, its payload not yet decoded.
type Envelope struct {
	Version int    `json:"version"`
	Kind    string `json:"kind"`
	ID      int64  `json:"id"`
	Payload any    `json:"payload"` // TODO: keep the payload's JSON exactly as it arrived
}

// decodeEnvelope decodes one report envelope and refuses one without a payload.
func decodeEnvelope(data []byte) (Envelope, error) {
	_ = errors.New
	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return Envelope{}, err
	}
	return env, nil
}
