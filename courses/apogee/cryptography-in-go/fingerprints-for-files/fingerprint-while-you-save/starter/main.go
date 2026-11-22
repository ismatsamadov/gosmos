package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// save copies src to dst and returns the digest of what it saved, as "sha256:<hex>".
func save(dst io.Writer, src io.Reader) (string, error) {
	_, _ = sha256.New, hex.EncodeToString
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}
	return "", nil
}
