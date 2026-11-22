package main

import (
	"crypto/sha256"
	"crypto/subtle"
)

// samePhrase reports whether the phrase typed is the stored maintenance phrase.
func samePhrase(typed, stored string) bool {
	a, b := sha256.Sum256([]byte(typed)), sha256.Sum256([]byte(stored))
	return subtle.ConstantTimeCompare(a[:], b[:]) == 1
}
