package main

import (
	"encoding/hex"
	"errors"
	"fmt"
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
