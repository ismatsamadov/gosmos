package main

import (
	"encoding/hex"
	"io"
)

// fingerprintAs returns the digest of everything r holds, named by its hash: "sha256:<hex>".
func fingerprintAs(r io.Reader, alg string) (string, error) {
	h, err := newHash(alg)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return alg + ":" + hex.EncodeToString(h.Sum(nil)), nil
}
