package main

import (
	"context"
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
	ID      string   // "furrowlight": lowercase, one word
	Name    string   // "Furrowlight"
	SatIDs  []string // the satellites it operates
	KeyHash string   // lowercase hex SHA-256 of its whole API key
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

// tokenMAC returns the HMAC-SHA256 of a token's claims under key.
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
	if err != nil {
		return "", unauthorized("invalid token")
	}
	if !now.Before(time.Unix(expires, 0)) {
		return "", unauthorized("token expired")
	}
	return parts[0], nil
}

// issueToken trades an API key for a token that lives one hour.
func (s *Service) issueToken(w http.ResponseWriter, r *http.Request) {
	key, err := credentials(r, "Key")
	c, ok := checkKey(s.Customers, key)
	if err == nil && !ok {
		err = unauthorized("invalid api key")
	}
	if err != nil {
		w.Header().Set("WWW-Authenticate", "Key")
		writeError(w, r, err)
		return
	}
	expires := s.Now().Add(time.Hour).Truncate(time.Second)
	tok := signToken(c.ID, expires, s.TokenKey)
	writeJSON(w, http.StatusOK, TokenResponse{Token: tok, ExpiresAt: expires})
}

// customerKey is the type of the context key a request's customer is kept under.
type customerKey struct{}

// withCustomer returns a copy of ctx that carries c.
func withCustomer(ctx context.Context, c Customer) context.Context {
	return context.WithValue(ctx, customerKey{}, c)
}

// customerFrom returns the customer ctx carries, if authenticate put one there.
func customerFrom(ctx context.Context) (Customer, bool) {
	c, ok := ctx.Value(customerKey{}).(Customer)
	return c, ok
}

// authenticate lets a request through only with a valid Bearer token, and passes its customer on in the request's context.
func authenticate(customers map[string]Customer, key []byte, now func() time.Time) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tok, err := credentials(r, "Bearer")
			id := ""
			if err == nil {
				id, err = verifyToken(tok, key, now())
			}
			c, known := customers[id]
			if err == nil && !known {
				err = unauthorized("invalid token")
			}
			if err != nil {
				w.Header().Set("WWW-Authenticate", "Bearer")
				writeError(w, r, err)
				return
			}
			next.ServeHTTP(w, r.WithContext(withCustomer(r.Context(), c)))
		})
	}
}

// private guards a customer's route: only a request with a valid token reaches h.
func (s *Service) private(h http.HandlerFunc) http.Handler {
	return authenticate(s.Customers, s.TokenKey, s.Now)(h)
}
