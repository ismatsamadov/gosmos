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
	if err := checkSig(pub, e.Receipt, e.Sig); err != nil {
		return Receipt{}, err
	}
	var r Receipt
	if err := decodeStrict(e.Receipt, &r); err != nil {
		return Receipt{}, fmt.Errorf("%w: %v", ErrBadEnvelope, err)
	}
	return r, nil
}
