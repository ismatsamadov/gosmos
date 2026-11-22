package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

// Signer signs canonical commands with uplinkd's key.
type Signer struct {
	key []byte
	mac hash.Hash // made once, for speed
}

// NewSigner returns a Signer for key.
func NewSigner(key []byte) *Signer {
	return &Signer{key: key, mac: hmac.New(sha256.New, key)}
}

// Sign returns the signature of a canonical command: its HMAC-SHA256, in hex.
func (s *Signer) Sign(canonical string) string {
	s.mac.Write([]byte(canonical))
	return hex.EncodeToString(s.mac.Sum(nil))
}
