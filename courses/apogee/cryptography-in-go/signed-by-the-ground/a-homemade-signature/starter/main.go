package main

import (
	"crypto/hmac"
	"crypto/sha256"
)

// tag returns the proof, sent with a command, that canonical was signed with key.
func tag(key []byte, canonical string) []byte {
	_ = hmac.New
	sum := sha256.Sum256(append(append([]byte{}, key...), canonical...))
	return sum[:]
}
