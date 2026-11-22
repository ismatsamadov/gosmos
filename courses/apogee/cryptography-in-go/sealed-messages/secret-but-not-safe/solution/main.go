package main

import (
	"crypto/aes"
	"crypto/cipher"
)

// seal locks plaintext under a 32-byte key and a 12-byte nonce and returns the locked bytes,
// its 16-byte tag at the end.
func seal(key, nonce, plaintext []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}
	return aead.Seal(nil, nonce, plaintext, nil)
}

// open unlocks what seal locked under the same key and nonce, or returns an error if any byte
// of it changed.
func open(key, nonce, sealed []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aead.Open(nil, nonce, sealed, nil)
}
