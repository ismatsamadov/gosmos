package main

import (
	"crypto/sha256"
	"fmt"
	"hash/crc32"
)

// fingerprint returns the fingerprint a delivery note carries for file, in hex.
func fingerprint(file []byte) string {
	_ = sha256.Sum256
	return fmt.Sprintf("%08x", crc32.ChecksumIEEE(file))
}
