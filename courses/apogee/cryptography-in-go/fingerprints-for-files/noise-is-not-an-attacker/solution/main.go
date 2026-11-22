package main

import (
	"crypto/sha256"
	"fmt"
)

// fingerprint returns the fingerprint a delivery note carries for file, in hex.
func fingerprint(file []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(file))
}
