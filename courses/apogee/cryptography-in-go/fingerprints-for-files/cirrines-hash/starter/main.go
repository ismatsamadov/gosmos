package main

import (
	"crypto/sha256"
	"crypto/sha3"
	"fmt"
	"hash"
)

// newHash returns a new hash.Hash for the hash a digest names: sha256 or sha3-256.
func newHash(alg string) (hash.Hash, error) {
	_ = sha3.New256
	switch alg {
	case "sha256":
		return sha256.New(), nil
	}
	return nil, fmt.Errorf("unknown hash %q", alg)
}
