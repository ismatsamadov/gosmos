package main

import (
	"crypto/hkdf"
	"crypto/sha256"
	"fmt"
	"slices"
)

// salt is HKDF's salt for every uplink key. It isn't secret, and production uses it too.
var salt = []byte("apogee uplink keys v1")

// fleet lists the satellites that hold keys of their own.
var fleet = []string{"pipit-1", "pipit-2", "pipit-3", "furrow-a", "furrow-b", "brine-3", "cirrine-1"}

// satKey returns satID's key for one purpose: "sign" checks its signed commands, "seal" opens its locked ones.
func satKey(master Key, purpose, satID string) ([]byte, error) {
	if purpose != "sign" && purpose != "seal" {
		return nil, fmt.Errorf("unknown purpose %q", purpose)
	}
	if !slices.Contains(fleet, satID) {
		return nil, fmt.Errorf("unknown satellite %q", satID)
	}
	return hkdf.Key(sha256.New, master.Secret, salt, purpose+"|"+satID, 32)
}
