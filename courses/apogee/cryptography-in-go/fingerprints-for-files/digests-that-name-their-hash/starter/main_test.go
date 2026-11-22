package main

import (
	"encoding/hex"
	"strings"
	"testing"
)

// The digests of "Apogee ground calling Pipit-1", as Apogee and Cirrine write them.
const (
	pipit1SHA256 = "sha256:fe5e0f73525ffc9f862167a80cb444fbc55de2a916041901a86964737b2fbd66"
	pipit1SHA3   = "sha3-256:7f7e6c3e2e3a3f1120a59350ba37cf6b382ad4644e840a4e15a35a179d9584ea"
)

// short shortens a digest longer than 50 characters to its hash's name and first 8 hex digits.
func short(s string) string {
	if len(s) <= 50 {
		return s
	}
	if name, digits, ok := strings.Cut(s, ":"); ok {
		return name + ":" + short(digits)
	}
	return s[:8] + "…"
}

func TestParseDigest(t *testing.T) {
	for _, s := range []string{pipit1SHA256, pipit1SHA3} {
		wantAlg, wantHex, _ := strings.Cut(s, ":")
		alg, sum, err := parseDigest(s)
		if err != nil || alg != wantAlg || hex.EncodeToString(sum) != wantHex {
			t.Errorf("parseDigest(%q) = %q, %d bytes, %v, want %q, 32 bytes, nil", short(s), alg, len(sum), err, wantAlg)
		}
	}
}

func TestParseDigestRefuses(t *testing.T) {
	digits := strings.TrimPrefix(pipit1SHA256, "sha256:")
	tests := []struct{ name, s, want string }{
		{"a bare digest", digits, "digest must be <hash>:<hex>"},
		{"SHA256 in capitals", "SHA256:" + digits, `unknown hash "SHA256"`},
		{"an md5 digest", "md5:" + digits[:32], `unknown hash "md5"`},
		{"63 digits", "sha256:" + digits[:63], "digest must be 64 lowercase hex digits"},
		{"uppercase hex", "sha256:" + strings.ToUpper(digits), "digest must be 64 lowercase hex digits"},
		{"a trailing space", pipit1SHA256 + " ", "digest must be 64 lowercase hex digits"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseDigest(tt.s)
			if err == nil || err.Error() != tt.want {
				t.Errorf("parseDigest(%q) error = %v, want %#q", short(tt.s), err, tt.want)
			}
		})
	}
}
