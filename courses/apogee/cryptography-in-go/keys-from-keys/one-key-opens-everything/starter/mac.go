package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// tag returns the HMAC-SHA256 of canonical under key.
func tag(key []byte, canonical string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(canonical))
	return mac.Sum(nil)
}

// verify reports whether sigHex is the signature of canonical under key.
func verify(canonical, sigHex string, key []byte) bool {
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		return false
	}
	return hmac.Equal(sig, tag(key, canonical))
}
