package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// save copies src to dst and returns the digest of what it saved, as "sha256:<hex>".
func save(dst io.Writer, src io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(dst, h), src); err != nil {
		return "", err
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}
