package main

import "crypto/rand"

// newKey returns a key named id whose secret is 32 random bytes.
func newKey(id string) Key {
	secret := make([]byte, 32)
	rand.Read(secret)
	return Key{ID: id, Secret: secret}
}

// newNonce returns 12 new bytes to send beside one locked command.
func newNonce() []byte {
	nonce := make([]byte, 12)
	rand.Read(nonce)
	return nonce
}

// newSalt returns 16 random bytes to keep beside one password.
func newSalt() []byte {
	salt := make([]byte, 16)
	rand.Read(salt)
	return salt
}

// newSessionToken returns a new token for one login to the ops console.
func newSessionToken() string {
	return rand.Text()
}
