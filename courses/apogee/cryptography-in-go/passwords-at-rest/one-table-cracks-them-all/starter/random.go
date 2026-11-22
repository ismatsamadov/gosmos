package main

import "crypto/rand"

// newSalt returns 16 random bytes to keep beside one password.
func newSalt() []byte {
	salt := make([]byte, 16)
	rand.Read(salt)
	return salt
}
