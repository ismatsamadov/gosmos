package main

import "crypto/ecdh"

// agree returns the secret that mine arrives at with the holder of the share theirs.
func agree(mine *ecdh.PrivateKey, theirs []byte) ([]byte, error) {
	pub, err := ecdh.X25519().NewPublicKey(theirs)
	if err != nil {
		return nil, err
	}
	return mine.ECDH(pub)
}
