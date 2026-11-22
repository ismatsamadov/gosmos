package main

import (
	"crypto/pbkdf2"
	"crypto/sha256"
)

// iterations is how many rounds of work derive spends on every password.
const iterations = 600000

// derive returns the 32-byte key the console stores for pw with the given salt.
func derive(pw string, salt []byte) ([]byte, error) {
	return pbkdf2.Key(sha256.New, pw, salt, iterations, 32)
}
