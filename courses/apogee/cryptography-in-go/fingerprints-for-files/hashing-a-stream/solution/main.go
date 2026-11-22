package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// fingerprintStream returns the SHA-256 digest, in hex, of everything r holds.
func fingerprintStream(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
