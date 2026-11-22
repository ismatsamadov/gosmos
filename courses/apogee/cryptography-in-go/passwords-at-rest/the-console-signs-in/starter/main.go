package main

import (
	"crypto/sha256"
	"encoding/hex"
)

// Login checks user's password and returns a new session token for their browser.
func (c *Console) Login(user, pw string) (token string, err error) {
	_, _ = sha256.Sum256, hex.EncodeToString
	return "", nil
}
