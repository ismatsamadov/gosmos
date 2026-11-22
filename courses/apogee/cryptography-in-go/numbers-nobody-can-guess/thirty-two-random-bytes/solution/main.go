package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// newKey returns a key named id whose secret is 32 random bytes.
func newKey(id string) Key {
	secret := make([]byte, 32)
	rand.Read(secret)
	return Key{ID: id, Secret: secret}
}

// parseSecret reads a key's secret from the 64 hex digits the key file holds.
func parseSecret(text string) ([]byte, error) {
	if len(text) != 64 {
		return nil, fmt.Errorf("secret has %d hex digits, want 64", len(text))
	}
	return hex.DecodeString(text)
}
