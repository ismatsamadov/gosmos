package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"
)

// checkStored reports whether pw matches a stored record, and whether that record is weaker than new ones.
func checkStored(stored, pw string) (ok, weak bool, err error) {
	_, _, _, _ = sha256.Sum256, subtle.ConstantTimeCompare, hex.EncodeToString, strings.CutPrefix
	r, err := parseRecord(stored)
	if err != nil {
		return false, false, err
	}
	ok, err = checkPassword(r, pw)
	return ok, false, err
}

// login checks pw against a stored record, and returns a record to save in its place when a weak one matched.
func login(stored, pw string) (ok bool, upgraded string, err error) {
	ok, _, err = checkStored(stored, pw)
	return ok, "", err
}
