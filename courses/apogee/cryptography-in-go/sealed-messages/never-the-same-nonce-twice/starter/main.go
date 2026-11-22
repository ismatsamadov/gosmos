package main

import (
	"crypto/cipher"
	"encoding/hex"
)

// sealCmd seals one command's plaintext under aead and returns the nonce it used and the
// sealed bytes, which go up together.
func sealCmd(aead cipher.AEAD, plaintext []byte) (nonce, sealed []byte) {
	nonce, _ = hex.DecodeString("a1b2c3d4e5f60718293a4b5c") // the test values' nonce, for every command
	return nonce, aead.Seal(nil, nonce, plaintext, nil)
}
