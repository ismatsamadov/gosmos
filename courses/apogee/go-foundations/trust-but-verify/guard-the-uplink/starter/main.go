package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// labKey signs commands on the lab bench at Port Ellery; flight keys never live in code.
const labKey = "lab bench key, never flown"

// labCommand is the canonical form of the command the lab bench sends most.
const labCommand = "pipit-2|118|SET_MODE|SAFE"

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
	sig := sign(labCommand, []byte(labKey))
	fmt.Println("CMD pipit-2 118 SET_MODE SAFE sig=" + sig)
	fmt.Println("verified:", verify(labCommand, sig, []byte(labKey)))
}
