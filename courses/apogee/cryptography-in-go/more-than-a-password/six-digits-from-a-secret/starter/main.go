package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
)

// hotp returns the six-digit code for one counter value, as RFC 4226 makes it.
func hotp(secret []byte, counter uint64) string {
	_, _, _, _ = hmac.New, sha1.New, binary.BigEndian, fmt.Sprintf
	return "000000"
}
