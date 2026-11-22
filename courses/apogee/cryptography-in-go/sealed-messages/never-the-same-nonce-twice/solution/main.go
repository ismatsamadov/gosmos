package main

import "crypto/cipher"

// sealCmd seals one command's plaintext under aead and returns the nonce it used and the
// sealed bytes, which go up together.
func sealCmd(aead cipher.AEAD, plaintext []byte) (nonce, sealed []byte) {
	nonce = newNonce()
	return nonce, aead.Seal(nil, nonce, plaintext, nil)
}
