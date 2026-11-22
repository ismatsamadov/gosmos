package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

// tokenKey is the test token key: 32 bytes, like the service's real key, which never appears in a lesson.
const tokenKey = "bookingd test token key 32 bytes"

// mint signs a token the way the service does, worked out here so the test needn't trust signToken.
func mint(customer string, expires time.Time, key string) string {
	claims := customer + "." + strconv.FormatInt(expires.Unix(), 10)
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(claims))
	return claims + "." + hex.EncodeToString(mac.Sum(nil))
}

// describe writes what verifyToken returned: the customer, then nil or the APIError's status, code and message.
func describe(customer string, err error) string {
	if err == nil {
		return fmt.Sprintf("%q, nil", customer)
	}
	if e, ok := errors.AsType[*APIError](err); ok {
		return fmt.Sprintf("%q, %d %s %q", customer, e.Status, e.Code, e.Message)
	}
	return fmt.Sprintf("%q, %q (not an *APIError)", customer, err)
}

func TestVerifyToken(t *testing.T) {
	eleven := time.Date(2028, 1, 5, 11, 0, 0, 0, time.UTC)
	ten := eleven.Add(-time.Hour)
	tok := mint("brinewatch", eleven, tokenKey)
	sig := tok[strings.LastIndex(tok, ".")+1:]
	tests := []struct {
		name string
		tok  string
		now  time.Time
		want string // the customer, or the message of the 401
	}{
		{"Brinewatch's token at 10:00", tok, ten, "brinewatch"},
		{"one second before it expires", tok, eleven.Add(-time.Second), "brinewatch"},
		{"at the second it expires", tok, eleven, "token expired"},
		{"an hour after it expired", tok, eleven.Add(time.Hour), "token expired"},
		{"the customer edited to cirrine", "cirrine.1830682800." + sig, ten, "invalid token"},
		{"the expiry moved a day later", "brinewatch.1830769200." + sig, ten, "invalid token"},
		{"an edited token that has expired too", "brinewatch.1830600000." + sig, ten, "invalid token"},
		{"signed with another key", mint("brinewatch", eleven, "another 32-byte key, not the one"), ten, "invalid token"},
		{"two parts", "brinewatch.1830682800", ten, "invalid token"},
		{"four parts", tok + ".cirrine", ten, "invalid token"},
		{"a signature that isn't hex", "brinewatch.1830682800." + strings.Repeat("zz", 32), ten, "invalid token"},
		{"half a signature", tok[:len(tok)-32], ten, "invalid token"},
		{"an empty token", "", ten, "invalid token"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := describe(tt.want, nil)
			if !strings.HasPrefix(tt.want, "brinewatch") {
				want = describe("", &APIError{Status: 401, Code: "unauthorized", Message: tt.want})
			}
			got := describe(verifyToken(tt.tok, []byte(tokenKey), tt.now))
			if got != want {
				t.Errorf("verifyToken(%q, key, %s) = %s, want %s", tt.tok, tt.now.Format("15:04:05"), got, want)
			}
		})
	}
}
