package main

import "encoding/json"

// Confirmation is the van's answer for one booking, as Hollis's sheet reads it.
type Confirmation struct {
	ID        string `json:"id"`
	Confirmed bool   `json:"confirmed"` // false means not confirmed yet
	Note      string `json:"note,omitempty"`
}

// encodeConfirmation writes one confirmation as JSON, without an empty note.
func encodeConfirmation(c Confirmation) ([]byte, error) {
	return json.Marshal(c)
}
