package main

import (
	"io"
	"strings"
	"testing"
)

// echo is the first message ever sent up from Apogee's new engineer, 29 bytes.
const echo = "Apogee ground calling Pipit-1"

// countingReader counts the calls to its Read, so a test can see whether anything was read.
type countingReader struct {
	r     io.Reader
	reads int
}

// Read counts the call, then reads from r.
func (c *countingReader) Read(b []byte) (int, error) {
	c.reads++
	return c.r.Read(b)
}

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

func TestFingerprintAs(t *testing.T) {
	tests := []struct{ alg, want string }{
		{"sha3-256", "sha3-256:7f7e6c3e2e3a3f1120a59350ba37cf6b382ad4644e840a4e15a35a179d9584ea"},
		{"sha256", "sha256:fe5e0f73525ffc9f862167a80cb444fbc55de2a916041901a86964737b2fbd66"},
	}
	for _, tt := range tests {
		got, err := fingerprintAs(strings.NewReader(echo), tt.alg)
		if err != nil || got != tt.want {
			t.Errorf("fingerprintAs(echo, %q) = %q, %v, want %q, nil", tt.alg, short(got), err, short(tt.want))
		}
	}
}

func TestUnknownHash(t *testing.T) {
	for _, alg := range []string{"md5", "sha3", "SHA3-256"} {
		r := &countingReader{r: strings.NewReader(echo)}
		got, err := fingerprintAs(r, alg)
		if err == nil {
			t.Errorf("fingerprintAs(echo, %q) = %q, nil, want an error", alg, short(got))
		}
		if r.reads > 0 {
			t.Errorf("fingerprintAs(echo, %q) read the file %d times before refusing, want 0", alg, r.reads)
		}
	}
}
