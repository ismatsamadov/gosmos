package main

import "fmt"

// sign returns the hex HMAC-SHA256 of a canonical command under key.
func sign(canonical string, key []byte) string {
	return ""
}

// verify reports whether sigHex is the signature of canonical under key.
func verify(canonical, sigHex string, key []byte) bool {
	return false
}

func main() {
	key := []byte("lab bench key, never flown")
	canonical := "pipit-2|118|SET_MODE|SAFE"
	sig := sign(canonical, key)
	fmt.Println("CMD pipit-2 118 SET_MODE SAFE sig=" + sig)
	fmt.Println("verified:", verify(canonical, sig, key))
}
