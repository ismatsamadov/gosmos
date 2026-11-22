package main

import (
	"crypto/hmac"
	"crypto/sha256"
)

// tag returns the proof, sent with a command, that canonical was signed with key.
func tag(key []byte, canonical string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(canonical))
	return mac.Sum(nil)
}
