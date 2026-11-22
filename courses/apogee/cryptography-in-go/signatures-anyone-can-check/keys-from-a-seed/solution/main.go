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

	priv := ed25519.NewKeyFromSeed(seed[:])
	pub := priv.Public().(ed25519.PublicKey)
	sig := ed25519.Sign(priv, receipt)
	verifies := ed25519.Verify(pub, receipt, sig)
	other := bytes.Clone(receipt)
	other[0] ^= 1
	changed := ed25519.Verify(pub, other, sig)

	fmt.Printf("public key: %d bytes, %x\n", len(pub), pub[:8])
	fmt.Printf("signature: %d bytes\n", len(sig))
	fmt.Printf("verifies: %t\n", verifies)
	fmt.Printf("one byte changed: %t\n", changed)
}
