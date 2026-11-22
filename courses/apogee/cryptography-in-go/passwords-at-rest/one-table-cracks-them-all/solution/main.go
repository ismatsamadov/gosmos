package main

import (
	"crypto/sha256"
	"encoding/hex"
)

// hashPassword returns the record the ops console keeps for pw.
func hashPassword(pw string) string {
	salt := newSalt()
	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(pw))
	return "salted$" + hex.EncodeToString(salt) + "$" + hex.EncodeToString(h.Sum(nil))
}
