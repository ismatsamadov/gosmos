package main

import "crypto/sha256"

// iterations is how many rounds of work derive spends on every password.
const iterations = 1

// derive returns the 32-byte key the console stores for pw with the given salt.
func derive(pw string, salt []byte) ([]byte, error) {
	sum := sha256.Sum256([]byte(string(salt) + pw))
	return sum[:], nil
}
