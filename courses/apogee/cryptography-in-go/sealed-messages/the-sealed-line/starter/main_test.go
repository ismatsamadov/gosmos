package main

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

// pingLine is Cirrine-1's first sealed command, PING, as the bible's test vector seals it: with
// the fixed test nonce, which real seals never use.
const (
	pingLine  = "SEAL 2 cirrine-1 4471 k-2028-07 a1b2c3d4e5f60718293a4b5c c1f2a9798d38081d764bbcfb5e35e2598d2584a4"
	nonceHex  = "a1b2c3d4e5f60718293a4b5c"
	sealedHex = "c1f2a9798d38081d764bbcfb5e35e2598d2584a4"
)

// withField returns pingLine with field i (counted from 0) replaced by v.
func withField(i int, v string) string {
	f := strings.Split(pingLine, " ")
	f[i] = v
	return strings.Join(f, " ")
}

func TestParsePing(t *testing.T) {
	s, err := parseSealed(pingLine)
	if err != nil {
		t.Fatalf("parseSealed(the sealed PING) = %v, want nil", err)
	}
	if s.Version != 2 || s.SatID != "cirrine-1" || s.Seq != 4471 || s.KeyID != "k-2028-07" {
		t.Errorf("parseSealed(the sealed PING) header = %d %q %d %q, want 2 \"cirrine-1\" 4471 \"k-2028-07\"",
			s.Version, s.SatID, s.Seq, s.KeyID)
	}
	if got := hex.EncodeToString(s.Nonce); got != nonceHex {
		t.Errorf("parseSealed(the sealed PING).Nonce = %s, want %s", got, nonceHex)
	}
	if got := hex.EncodeToString(s.Ciphertext); got != sealedHex {
		t.Errorf("parseSealed(the sealed PING).Ciphertext = %s, want %s", got, sealedHex)
	}
}

func TestFormatPing(t *testing.T) {
	nonce, _ := hex.DecodeString(nonceHex)
	sealed, _ := hex.DecodeString(sealedHex)
	s := Sealed{Version: 2, SatID: "cirrine-1", Seq: 4471, KeyID: "k-2028-07", Nonce: nonce, Ciphertext: sealed}
	if got := formatSealed(s); got != pingLine {
		t.Errorf("formatSealed(the sealed PING) = %q, want %q", got, pingLine)
	}
	back, err := parseSealed(pingLine)
	if got := formatSealed(back); err != nil || got != pingLine {
		t.Errorf("formatSealed(parseSealed(the sealed PING)) = %q, %v, want the line byte for byte", got, err)
	}
}

func TestMalformed(t *testing.T) {
	tests := []struct{ name, line string }{
		{"6 fields", pingLine[:strings.LastIndex(pingLine, " ")]},
		{"8 fields", pingLine + " PING"},
		{"a doubled space", strings.Replace(pingLine, " 4471", "  4471", 1)},
		{"a space at the end", pingLine + " "},
		{"CMD in place of SEAL", withField(0, "CMD")},
		{"CMD and seven fields", "CMD cirrine-1 4472 SET_MODE SAFE NOW sig=" + strings.Repeat("ab", 32)},
		{"no satellite, two spaces", withField(2, "")},
		{"the satellite Cirrine-1", withField(2, "Cirrine-1")},
		{"sequence number 0", withField(3, "0")},
		{"sequence number 44a1", withField(3, "44a1")},
		{"a 22-digit nonce", withField(5, nonceHex[:22])},
		{"an odd-length ciphertext", withField(6, sealedHex[:39])},
		{"an uppercase ciphertext", withField(6, strings.ToUpper(sealedHex))},
		{"a 16-byte ciphertext", withField(6, strings.Repeat("ab", 16))},
		{"a 273-byte ciphertext", withField(6, strings.Repeat("ab", 273))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseSealed(tt.line); !errors.Is(err, ErrMalformed) {
				t.Errorf("parseSealed(a line with %s) = %v, want an error wrapping ErrMalformed", tt.name, err)
			}
		})
	}
}

func TestCiphertextSizes(t *testing.T) {
	for _, n := range []int{17, 272} {
		if _, err := parseSealed(withField(6, strings.Repeat("ab", n))); err != nil {
			t.Errorf("parseSealed(a line with a %d-byte ciphertext) = %v, want nil", n, err)
		}
	}
}

func TestVersion(t *testing.T) {
	if _, err := parseSealed(withField(1, "3")); !errors.Is(err, ErrVersion) {
		t.Errorf("parseSealed(SEAL 3 ...) = %v, want an error wrapping ErrVersion", err)
	}
}
