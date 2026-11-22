package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"
)

// Customer is one company that books passes through the service.
type Customer struct {
	ID      string // "furrowlight": lowercase, one word
	Name    string // "Furrowlight"
	KeyHash string // lowercase hex SHA-256 of its whole API key
}

// unauthorized returns the 401 that tells the caller msg.
func unauthorized(msg string) error {
	return &APIError{Status: http.StatusUnauthorized, Code: "unauthorized", Message: msg}
}

// credentials returns the credential that follows scheme in r's Authorization header.
func credentials(r *http.Request, scheme string) (string, error) {
	got, credential, _ := strings.Cut(r.Header.Get("Authorization"), " ")
	if !strings.EqualFold(got, scheme) || credential == "" {
		return "", unauthorized("missing credentials")
	}
	return credential, nil
}

// checkKey returns the customer whose API key this is.
func checkKey(customers map[string]Customer, key string) (Customer, bool) {
	_ = sha256.Sum256
	_ = hex.EncodeToString
	id, _, _ := strings.Cut(key, ".")
	c, ok := customers[id]
	// TODO: compare the key's SHA-256, in hex, with KeyHash instead of the key itself.
	if !ok || subtle.ConstantTimeCompare([]byte(key), []byte(c.KeyHash)) != 1 {
		return Customer{}, false
	}
	return c, true
}
