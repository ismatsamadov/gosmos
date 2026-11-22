package main

import (
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

// parseShare reads a station's public share from the 64 lowercase hex digits it sends.
func parseShare(text string) (*ecdh.PublicKey, error) {
	b, err := hex.DecodeString(text)
	if err != nil || len(b) != 32 || hex.EncodeToString(b) != text {
		return nil, errors.New("share must be 64 lowercase hex digits")
	}
	return ecdh.X25519().NewPublicKey(b)
}

// sessionKey returns the key that seals one backhaul session, made from the secret both ends agreed on.
func sessionKey(shared, dialerShare, answererShare []byte) ([]byte, error) {
	if len(shared) != 32 {
		return nil, fmt.Errorf("shared secret is %d bytes, want 32", len(shared))
	}
	info := "apogee backhaul v1|" + hex.EncodeToString(dialerShare) + "|" + hex.EncodeToString(answererShare)
	return hkdf.Key(sha256.New, shared, nil, info, 32)
}
