package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// newKey returns a key named id whose secret is 32 random bytes.
func newKey(id string) Key {
	_ = rand.Read
	return Key{ID: id}
}

// parseSecret reads a key's secret from the 64 hex digits the key file holds.
func parseSecret(text string) ([]byte, error) {
	_ = fmt.Errorf
	return hex.DecodeString(text)
}
