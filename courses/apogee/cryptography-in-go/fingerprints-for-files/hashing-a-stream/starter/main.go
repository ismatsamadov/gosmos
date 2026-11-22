package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// fingerprintStream returns the SHA-256 digest, in hex, of everything r holds.
func fingerprintStream(r io.Reader) (string, error) {
	_ = io.Copy
	buf := make([]byte, 64*1024)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	sum := sha256.Sum256(buf[:n])
	return hex.EncodeToString(sum[:]), nil
}
