package main

import (
	"crypto/aes"
	"crypto/cipher"
)

// seal locks plaintext under a 32-byte key and a 12-byte nonce and returns the locked bytes.
func seal(key, nonce, plaintext []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	sealed := make([]byte, len(plaintext))
	cipher.NewCTR(block, counter(nonce)).XORKeyStream(sealed, plaintext)
	return sealed
}

// open unlocks what seal locked under the same key and nonce.
func open(key, nonce, sealed []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(sealed))
	cipher.NewCTR(block, counter(nonce)).XORKeyStream(plaintext, sealed)
	return plaintext, nil
}

// counter pads a 12-byte nonce to the 16 bytes counter mode starts from.
func counter(nonce []byte) []byte {
	iv := make([]byte, aes.BlockSize)
	copy(iv, nonce)
	return iv
}
