package main

import (
	"crypto/hkdf"
	"crypto/sha256"
	"errors"
)

// salt is HKDF's salt for every uplink key. It isn't secret, and production uses it too.
var salt = []byte("apogee uplink keys v1")

// signKey returns the key satID's flight software checks signed commands with.
func signKey(master Key, satID string) ([]byte, error) {
	if satID == "" {
		return nil, errors.New("no satellite ID")
	}
	return hkdf.Key(sha256.New, master.Secret, salt, "sign|"+satID, 32)
}
