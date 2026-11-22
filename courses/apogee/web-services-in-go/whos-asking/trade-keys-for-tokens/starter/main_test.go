package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// The beta customers' test API keys, as Yusra gave them out. The service keeps none of them.
const (
	furrowlightKey = "furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b"
	cirrineKey     = "cirrine.c3a1f7e95d2b4086a7e1c9f3b5d8a024"
)

// tokenKey is the test token key: 32 bytes, like the service's real key, which never appears in a lesson.
const tokenKey = "bookingd test token key 32 bytes"

// customers returns the customer table as the service keeps it: a hash of each key, never the key.
func customers() map[string]Customer {
	return map[string]Customer{
		"apogee":      {ID: "apogee", Name: "Apogee", KeyHash: "8018b533030c734641dda97af10c84dd57357d0e37b7a0bdee7ed2a63c911407"},
		"furrowlight": {ID: "furrowlight", Name: "Furrowlight", KeyHash: "95ff036b9ad2e59c26fae1be7dda78b12412d243fcbb6567bc8124e6f1c93fd3"},
		"brinewatch":  {ID: "brinewatch", Name: "Brinewatch", KeyHash: "6f5748f2886bd91e86930360f00a38ae260f6de18ea20012f4e4ec718d483f85"},
		"cirrine":     {ID: "cirrine", Name: "Cirrine Atmospherics", KeyHash: "d9ee2e5de1796bbeba223c0737933b3e60d53687d479f8a04c98c99df87b510d"},
	}
}

// mint signs a token the way the service does, worked out here so the test needn't trust signToken.
func mint(customer string, expires time.Time, key string) string {
	claims := customer + "." + strconv.FormatInt(expires.Unix(), 10)
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(claims))
	return claims + "." + hex.EncodeToString(mac.Sum(nil))
}

// postTokens sends POST /tokens with auth as its Authorization header ("" sends none) to a service whose clock reads now.
func postTokens(now time.Time, auth string) *httptest.ResponseRecorder {
	s := &Service{Customers: customers(), TokenKey: []byte(tokenKey), Now: func() time.Time { return now }}
	req := httptest.NewRequest(http.MethodPost, "/tokens", nil)
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	rec := httptest.NewRecorder()
	http.HandlerFunc(s.issueToken).ServeHTTP(rec, req)
	return rec
}

func TestIssueToken(t *testing.T) {
	now := time.Date(2028, 1, 5, 10, 0, 0, 600_000_000, time.UTC) // Wednesday, 10:00:00.6
	eleven := time.Date(2028, 1, 5, 11, 0, 0, 0, time.UTC)
	for _, tc := range []struct{ name, customer, key string }{
		{"Furrowlight", "furrowlight", furrowlightKey},
		{"Cirrine", "cirrine", cirrineKey},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := postTokens(now, "Key "+tc.key)
			if rec.Code != http.StatusOK {
				t.Fatalf("POST /tokens with %s's key = %d %s, want 200", tc.name, rec.Code, bytes.TrimSpace(rec.Body.Bytes()))
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("POST /tokens Content-Type = %q, want %q", ct, "application/json")
			}
			var got TokenResponse
			dec := json.NewDecoder(rec.Body)
			dec.DisallowUnknownFields()
			if err := dec.Decode(&got); err != nil {
				t.Fatalf("decoding the TokenResponse of POST /tokens: %v", err)
			}
			if !got.ExpiresAt.Equal(eleven) {
				t.Errorf("POST /tokens at 10:00:00.6 expires_at = %s, want %s", got.ExpiresAt.Format(time.RFC3339Nano), eleven.Format(time.RFC3339Nano))
			}
			if want := mint(tc.customer, eleven, tokenKey); got.Token != want {
				t.Errorf("POST /tokens at 10:00:00.6 token = %q, want %q", got.Token, want)
			}
		})
	}
}

func TestIssueTokenRefuses(t *testing.T) {
	now := time.Date(2028, 1, 5, 10, 0, 0, 600_000_000, time.UTC)
	tests := []struct {
		name string
		auth string // "" sends no Authorization header
		want string // the message of the 401
	}{
		{"no Authorization header", "", "missing credentials"},
		{"the key sent as a Bearer token", "Bearer " + furrowlightKey, "missing credentials"},
		{"a wrong secret", "Key furrowlight.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6c", "invalid api key"},
		{"a customer the service doesn't know", "Key kessick.5b0e7c1d9a3f4e6b8c2d1a0f9e8d7c6b", "invalid api key"},
		{"a token sent as the key", "Key " + mint("furrowlight", now.Add(time.Hour), tokenKey), "invalid api key"},
		{"Furrowlight's KeyHash sent as the key", "Key " + customers()["furrowlight"].KeyHash, "invalid api key"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := postTokens(now, tt.auth)
			if rec.Code != http.StatusUnauthorized {
				t.Errorf("POST /tokens with %s = %d, want 401", tt.name, rec.Code)
			}
			if got := rec.Header().Get("WWW-Authenticate"); got != "Key" {
				t.Errorf("POST /tokens with %s WWW-Authenticate = %q, want %q", tt.name, got, "Key")
			}
			var got ErrorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("POST /tokens with %s body = %q, want an ErrorBody: %v", tt.name, rec.Body, err)
			}
			if want := (ErrorBody{Code: "unauthorized", Message: tt.want}); got != want {
				t.Errorf("POST /tokens with %s = %+v, want %+v", tt.name, got, want)
			}
		})
	}
}
