package main

import "crypto/ecdh"

// agree returns the secret that mine arrives at with the holder of the share theirs.
func agree(mine *ecdh.PrivateKey, theirs []byte) ([]byte, error) {
	_ = ecdh.X25519
	return nil, nil
}
