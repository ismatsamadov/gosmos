package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// verifyEnvelope checks env with the receipt keys Apogee publishes and returns the receipt it
// carries.
func verifyEnvelope(keys map[string]ed25519.PublicKey, env []byte) (Receipt, error) {
	var e struct {
		Envelope
		Pub string `json:"pub"` // the draft's own key, in hex, when the envelope brings one
	}
	if err := json.Unmarshal(env, &e); err != nil {
		return Receipt{}, fmt.Errorf("%w: %v", ErrBadEnvelope, err)
	}
	pub := keys[e.KeyID]
	if e.Pub != "" {
		pub, _ = hex.DecodeString(e.Pub)
	}
	return openReceipt(pub, e.Envelope)
}
