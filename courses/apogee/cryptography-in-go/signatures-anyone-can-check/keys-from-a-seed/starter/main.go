package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
)

func main() {
	seed := sha256.Sum256([]byte("receipts r-2028-08")) // r-2028-08's test seed, readable on purpose
	receipt := []byte("bk-5274 furrow-a-20280827-1402.tif landed 2028-08-27T14:13:05Z")

	// TODO: make the key from seed, sign receipt, and check the signature twice.
	var pub ed25519.PublicKey
	var sig []byte
	verifies, changed := false, false
	_, _, _ = seed, receipt, bytes.Clone

	fmt.Printf("public key: %d bytes, %x\n", len(pub), pub)
	fmt.Printf("signature: %d bytes\n", len(sig))
	fmt.Printf("verifies: %t\n", verifies)
	fmt.Printf("one byte changed: %t\n", changed)
}
