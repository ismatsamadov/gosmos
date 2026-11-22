package main

import (
	"crypto/ecdh"
	"encoding/hex"
	"errors"
)

// parseShare reads a station's public share from the 64 lowercase hex digits it sends.
func parseShare(text string) (*ecdh.PublicKey, error) {
	b, err := hex.DecodeString(text)
	if err != nil || len(b) != 32 || hex.EncodeToString(b) != text {
		return nil, errors.New("share must be 64 lowercase hex digits")
	}
	return ecdh.X25519().NewPublicKey(b)
}
