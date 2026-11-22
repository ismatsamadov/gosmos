package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"slices"
	"strings"
)

// Customer is one company that books passes through the service.
type Customer struct {
	ID         string   // "furrowlight": lowercase, one word
	Name       string   // "Furrowlight"
	SatIDs     []string // the satellites it operates
	KeyHash    string   // lowercase hex SHA-256 of its whole API key
	WebhookURL string   // where its events go
	WebhookKey []byte   // the secret its events are signed with
}

// checkKey returns the customer whose API key this is.
func checkKey(customers map[string]Customer, key string) (Customer, bool) {
	id, _, _ := strings.Cut(key, ".")
	c, ok := customers[id]
	sum := sha256.Sum256([]byte(key))
	hash := hex.EncodeToString(sum[:])
	if !ok || subtle.ConstantTimeCompare([]byte(hash), []byte(c.KeyHash)) != 1 {
		return Customer{}, false
	}
	return c, true
}

// owns reports whether c operates the satellite satID.
func owns(c Customer, satID string) bool {
	return slices.Contains(c.SatIDs, satID)
}
