package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// hmacOf returns the HMAC-SHA256 of a canonical command under key.
func hmacOf(canonical string, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(canonical))
	return mac.Sum(nil)
}

// sign returns the hex HMAC-SHA256 of a canonical command under key.
func sign(canonical string, key []byte) string {
	return hex.EncodeToString(hmacOf(canonical, key))
}

// verify reports whether sigHex is the signature of canonical under key.
func verify(canonical, sigHex string, key []byte) bool {
	sig, err := hex.DecodeString(sigHex)
	if err != nil {
		return false
	}
	return hmac.Equal(sig, hmacOf(canonical, key))
}

func main() {
	key := []byte("lab bench key, never flown")
	canonical := "pipit-2|118|SET_MODE|SAFE"
	sig := sign(canonical, key)
	fmt.Println("CMD pipit-2 118 SET_MODE SAFE sig=" + sig)
	fmt.Println("verified:", verify(canonical, sig, key))
}
