package main

import (
	"crypto/sha256"
	"crypto/subtle"
)

// samePhrase reports whether the phrase typed is the stored maintenance phrase.
func samePhrase(typed, stored string) bool {
	_, _ = sha256.Sum256, subtle.ConstantTimeCompare
	return typed == stored
}
