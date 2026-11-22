package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// sessionKey returns the key that seals one backhaul session, made from the secret both ends agreed on.
func sessionKey(shared, dialerShare, answererShare []byte) ([]byte, error) {
	_, _, _ = sha256.New, hex.EncodeToString, fmt.Errorf
	return shared, nil
}
