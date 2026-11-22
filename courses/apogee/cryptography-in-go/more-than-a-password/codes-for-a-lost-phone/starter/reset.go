package main

import (
	"crypto/sha256"
	"encoding/hex"
)

// tokenDigest returns the hex SHA-256 of a token: what the console keeps instead of it.
func tokenDigest(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
