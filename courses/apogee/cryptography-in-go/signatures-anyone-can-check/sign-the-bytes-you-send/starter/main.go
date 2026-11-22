package main

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
)

// openEnvelope checks env's signature with pub and returns the receipt it carries.
func openEnvelope(pub ed25519.PublicKey, env []byte) (Receipt, error) {
	var e Envelope
	if err := json.Unmarshal(env, &e); err != nil {
		return Receipt{}, fmt.Errorf("%w: %v", ErrBadEnvelope, err)
	}
	var r Receipt
	if err := json.Unmarshal(e.Receipt, &r); err != nil {
		return Receipt{}, fmt.Errorf("%w: %v", ErrBadEnvelope, err)
	}
	again, err := json.Marshal(r) // the receipt as Apogee writes it
	if err != nil {
		return Receipt{}, err
	}
	if err := checkSig(pub, again, e.Sig); err != nil {
		return Receipt{}, err
	}
	return r, nil
}
