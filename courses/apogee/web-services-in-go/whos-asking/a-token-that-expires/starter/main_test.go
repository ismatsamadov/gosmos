package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"
)

// tokenKey is the test token key: 32 bytes, like the service's real key, which never appears in a lesson.
const tokenKey = "bookingd test token key 32 bytes"

// hmacHex returns the hex HMAC-SHA256 of claims under key, worked out here so the test needn't trust signToken.
func hmacHex(claims, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(claims))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestSignToken(t *testing.T) {
	eleven := time.Date(2028, 1, 5, 11, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		customer string
		expires  time.Time
		key      string
		want     string
	}{
		{"Furrowlight's token until 11:00", "furrowlight", eleven, tokenKey,
			"furrowlight.1830682800." + hmacHex("furrowlight.1830682800", tokenKey)},
		{"an hour later: a new signature, not only a new expiry", "furrowlight", eleven.Add(time.Hour), tokenKey,
			"furrowlight.1830686400." + hmacHex("furrowlight.1830686400", tokenKey)},
		{"Brinewatch's token until 11:00", "brinewatch", eleven, tokenKey,
			"brinewatch.1830682800." + hmacHex("brinewatch.1830682800", tokenKey)},
		{"another key: another signature", "furrowlight", eleven, "another 32-byte key, not the one",
			"furrowlight.1830682800." + hmacHex("furrowlight.1830682800", "another 32-byte key, not the one")},
		{"Hollis's 12:00 OUR TIME is 11:00 UTC", "furrowlight", time.Date(2028, 1, 5, 12, 0, 0, 0, time.FixedZone("", 60*60)), tokenKey,
			"furrowlight.1830682800." + hmacHex("furrowlight.1830682800", tokenKey)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := signToken(tt.customer, tt.expires, []byte(tt.key)); got != tt.want {
				t.Errorf("signToken(%q, %s, %q) = %q, want %q", tt.customer, tt.expires.Format(time.RFC3339), tt.key, got, tt.want)
			}
		})
	}
}
