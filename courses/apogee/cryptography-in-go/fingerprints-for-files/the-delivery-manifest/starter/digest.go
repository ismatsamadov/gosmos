package main

import (
	"crypto/sha256"
	"crypto/sha3"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"strings"
)

// parseDigest reads a digest written as <hash>:<hex>, such as "sha256:fe5e…".
func parseDigest(s string) (alg string, sum []byte, err error) {
	alg, hexText, ok := strings.Cut(s, ":")
	if !ok {
		return "", nil, errors.New("digest must be <hash>:<hex>")
	}
	if alg != "sha256" && alg != "sha3-256" {
		return "", nil, fmt.Errorf("unknown hash %q", alg)
	}
	if len(hexText) != 64 || strings.Trim(hexText, "0123456789abcdef") != "" {
		return "", nil, errors.New("digest must be 64 lowercase hex digits")
	}
	sum, err = hex.DecodeString(hexText)
	return alg, sum, err
}

// newHash returns a new hash.Hash for the hash a digest names: sha256 or sha3-256.
func newHash(alg string) (hash.Hash, error) {
	switch alg {
	case "sha256":
		return sha256.New(), nil
	case "sha3-256":
		return sha3.New256(), nil
	}
	return nil, fmt.Errorf("unknown hash %q", alg)
}

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
