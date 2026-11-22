package main

import (
	"crypto/ecdh"
	"encoding/hex"
	"errors"
)

// parseShare reads a station's public share from the 64 lowercase hex digits it sends.
func parseShare(text string) (*ecdh.PublicKey, error) {
	_ = errors.New
	b, _ := hex.DecodeString(text)
	return ecdh.X25519().NewPublicKey(b)
}

// agreeHex returns the secret that mine agrees on with the station whose share is text.
func agreeHex(mine *ecdh.PrivateKey, text string) ([]byte, error) {
	theirs, err := parseShare(text)
	if err != nil {
		return nil, err
	}
	secret, _ := mine.ECDH(theirs)
	return secret, nil
}
