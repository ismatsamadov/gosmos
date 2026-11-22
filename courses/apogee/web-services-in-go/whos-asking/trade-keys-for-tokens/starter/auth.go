package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Customer is one company that books passes through the service.
type Customer struct {
	ID      string // "furrowlight": lowercase, one word
	Name    string // "Furrowlight"
	KeyHash string // lowercase hex SHA-256 of its whole API key
}

// TokenResponse is the body of POST /tokens.
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
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
	id, _, _ := strings.Cut(key, ".")
	c, ok := customers[id]
	sum := sha256.Sum256([]byte(key))
	hash := hex.EncodeToString(sum[:])
	if !ok || subtle.ConstantTimeCompare([]byte(hash), []byte(c.KeyHash)) != 1 {
		return Customer{}, false
	}
	return c, true
}

// tokenMAC returns the HMAC-SHA256 of a token's claims, "<customer>.<expiry>", under key.
func tokenMAC(claims string, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(claims))
	return mac.Sum(nil)
}

// signToken returns a token that names customer until expires, signed with key.
func signToken(customer string, expires time.Time, key []byte) string {
	claims := customer + "." + strconv.FormatInt(expires.Unix(), 10)
	return claims + "." + hex.EncodeToString(tokenMAC(claims, key))
}

// verifyToken returns the customer a token names, if its signature checks out and it hasn't expired at now.
func verifyToken(tok string, key []byte, now time.Time) (string, error) {
	parts := strings.Split(tok, ".")
	if len(parts) != 3 {
		return "", unauthorized("invalid token")
	}
	sig, err := hex.DecodeString(parts[2])
	if err != nil || !hmac.Equal(sig, tokenMAC(parts[0]+"."+parts[1], key)) {
		return "", unauthorized("invalid token")
	}
	expires, err := strconv.ParseInt(parts[1], 10, 64)
	switch {
	case err != nil:
		return "", unauthorized("invalid token")
	case !now.Before(time.Unix(expires, 0)):
		return "", unauthorized("token expired")
	}
	return parts[0], nil
}

// issueToken trades an API key for a token that lives one hour.
func (s *Service) issueToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
