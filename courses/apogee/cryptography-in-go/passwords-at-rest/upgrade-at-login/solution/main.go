package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"
)

// checkStored reports whether pw matches a stored record, and whether that record is weaker than new ones.
func checkStored(stored, pw string) (ok, weak bool, err error) {
	if old, found := strings.CutPrefix(stored, "sha256$"); found {
		sum := sha256.Sum256([]byte(pw))
		typed := hex.EncodeToString(sum[:])
		return subtle.ConstantTimeCompare([]byte(typed), []byte(old)) == 1, true, nil
	}
	r, err := parseRecord(stored)
	if err != nil {
		return false, false, err
	}
	ok, err = checkPassword(r, pw)
	return ok, r.Iterations < iterations, err
}

// login checks pw against a stored record, and returns a record to save in its place when a weak one matched.
func login(stored, pw string) (ok bool, upgraded string, err error) {
	ok, weak, err := checkStored(stored, pw)
	if err != nil || !ok || !weak {
		return ok, "", err
	}
	upgraded, err = newRecord(pw)
	return err == nil, upgraded, err
}
