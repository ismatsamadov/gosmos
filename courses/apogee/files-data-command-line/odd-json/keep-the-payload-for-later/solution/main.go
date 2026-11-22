package main

import (
	"encoding/json"
	"errors"
)

// Envelope is one Furrowlight delivery report, its payload not yet decoded.
type Envelope struct {
	Version int             `json:"version"`
	Kind    string          `json:"kind"`
	ID      int64           `json:"id"`
	Payload json.RawMessage `json:"payload"`
}

// decodeEnvelope decodes one report envelope and refuses one without a payload.
func decodeEnvelope(data []byte) (Envelope, error) {
	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return Envelope{}, err
	}
	if len(env.Payload) == 0 || string(env.Payload) == "null" {
		return Envelope{}, errors.New("report has no payload")
	}
	return env, nil
}
