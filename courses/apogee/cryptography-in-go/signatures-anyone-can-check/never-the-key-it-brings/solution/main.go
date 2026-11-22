package main

import (
	"crypto/ed25519"
	"fmt"
)

// verifyEnvelope checks env with the receipt keys Apogee publishes and returns the receipt it
// carries.
func verifyEnvelope(keys map[string]ed25519.PublicKey, env []byte) (Receipt, error) {
	e, err := decodeEnvelope(env)
	if err != nil {
		return Receipt{}, err
	}
	pub, ok := keys[e.KeyID]
	if !ok {
		return Receipt{}, fmt.Errorf("%w: %q", ErrUnknownSigner, e.KeyID)
	}
	r, err := openReceipt(pub, e)
	if err != nil {
		return Receipt{}, err
	}
	if r.KeyID != e.KeyID {
		return Receipt{}, fmt.Errorf("%w: receipt names key %q", ErrBadEnvelope, r.KeyID)
	}
	return r, nil
}
