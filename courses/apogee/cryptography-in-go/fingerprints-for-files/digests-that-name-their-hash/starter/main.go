package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// parseDigest reads a digest written as <hash>:<hex>, such as "sha256:fe5e…".
func parseDigest(s string) (alg string, sum []byte, err error) {
	_, _ = errors.New, fmt.Errorf
	alg, hexText, _ := strings.Cut(s, ":")
	sum, err = hex.DecodeString(hexText)
	return alg, sum, err
}
