package main

import (
	"crypto/sha256"
	"crypto/sha3"
	"fmt"
	"hash"
)

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
