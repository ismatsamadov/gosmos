package main

import "crypto/hmac"

// valid reports whether t is the tag of canonical under key, as Pipit-3 checks it.
func valid(key []byte, canonical string, t []byte) bool {
	return hmac.Equal(t, tag(key, canonical))
}
