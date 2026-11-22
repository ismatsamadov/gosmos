package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Signer signs canonical commands with uplinkd's key.
type Signer struct {
	key []byte
}

// NewSigner returns a Signer for key.
func NewSigner(key []byte) *Signer {
	return &Signer{key: key}
}

// Sign returns the signature of a canonical command: its HMAC-SHA256, in hex.
func (s *Signer) Sign(canonical string) string {
	mac := hmac.New(sha256.New, s.key)
	mac.Write([]byte(canonical))
	return hex.EncodeToString(mac.Sum(nil))
}
