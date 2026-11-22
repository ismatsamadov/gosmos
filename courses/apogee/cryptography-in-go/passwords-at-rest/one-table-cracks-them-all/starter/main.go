package main

import (
	"crypto/sha256"
	"encoding/hex"
)

// hashPassword returns the record the ops console keeps for pw.
func hashPassword(pw string) string {
	sum := sha256.Sum256([]byte(pw))
	return "sha256$" + hex.EncodeToString(sum[:])
}
