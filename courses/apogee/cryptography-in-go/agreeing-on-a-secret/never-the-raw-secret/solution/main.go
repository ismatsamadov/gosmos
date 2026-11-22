package main

import (
	"crypto/hkdf"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// sessionKey returns the key that seals one backhaul session, made from the secret both ends agreed on.
func sessionKey(shared, dialerShare, answererShare []byte) ([]byte, error) {
	if len(shared) != 32 {
		return nil, fmt.Errorf("shared secret is %d bytes, want 32", len(shared))
	}
	info := "apogee backhaul v1|" + hex.EncodeToString(dialerShare) + "|" + hex.EncodeToString(answererShare)
	return hkdf.Key(sha256.New, shared, nil, info, 32)
}
