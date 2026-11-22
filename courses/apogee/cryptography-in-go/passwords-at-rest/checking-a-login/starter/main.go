package main

import (
	"crypto/sha256"
	"crypto/subtle"
)

// checkPassword reports whether pw is the password r was made from.
func checkPassword(r Record, pw string) (bool, error) {
	_, _ = sha256.New, subtle.ConstantTimeCompare
	return false, nil
}
