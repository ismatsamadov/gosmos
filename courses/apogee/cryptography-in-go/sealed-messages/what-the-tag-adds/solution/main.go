package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

func main() {
	// Cirrine-1's test seal key from k-2028-07, and the tests' fixed nonce: test values only.
	key, _ := hex.DecodeString("427a8ae78a295915acb18a87bdbb40be20d49d27c15dac44e320d4510e070171")
	nonce, _ := hex.DecodeString("a1b2c3d4e5f60718293a4b5c")
	block, _ := aes.NewCipher(key)
	aead, _ := cipher.NewGCM(block)

	plaintext := []byte("SET_MODE|SAFE")
	sealed := aead.Seal(nil, nonce, plaintext, nil)
	tampered := bytes.Clone(sealed)
	tampered[len(tampered)-1] ^= 1
	_, err := aead.Open(nil, nonce, tampered, nil)

	fmt.Printf("plaintext: %d bytes\n", len(plaintext))
	fmt.Printf("sealed: %d bytes\n", len(sealed))
	fmt.Printf("nonce: %d bytes, tag: %d bytes\n", aead.NonceSize(), aead.Overhead())
	fmt.Printf("tampered: %v\n", err)
}
