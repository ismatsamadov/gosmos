package main

import (
	"bytes"
	"crypto/ecdh"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

// The shares of the test keys "ephemeral PEL" and "ephemeral YST", and the secret they agree on.
const (
	pelShare = "a10c608ed294ad1d437b121ce63f146ee081f200b8b71126bc8f5642e69aa93e"
	ystShare = "b6f8a5a974e292828807a52f32334ac3925ea872ca82de318f53e971fafc5a66"
	agreed   = "c50d64c87ddd744d6c24842274ac98d74c33be2bdadd89463eae753cae91cc7a"
)

// pelKey returns Port Ellery's test X25519 key, whose private bytes are the SHA-256 of "ephemeral PEL".
func pelKey(t *testing.T) *ecdh.PrivateKey {
	t.Helper()
	seed := sha256.Sum256([]byte("ephemeral PEL"))
	key, err := ecdh.X25519().NewPrivateKey(seed[:])
	if err != nil {
		t.Fatalf("ecdh.X25519().NewPrivateKey(SHA-256 of %q) = %v", "ephemeral PEL", err)
	}
	return key
}

// short returns s, or its first 8 characters and an ellipsis when it is longer than 50.
func short(s string) string {
	if len(s) > 50 {
		return s[:8] + "…"
	}
	return s
}

func TestParseShare(t *testing.T) {
	share, err := parseShare(pelShare)
	if err != nil {
		t.Fatalf("parseShare(%q) = %v, want PEL's share", short(pelShare), err)
	}
	if got := hex.EncodeToString(share.Bytes()); got != pelShare {
		t.Errorf("parseShare(%q).Bytes() = %s, want %s", short(pelShare), short(got), short(pelShare))
	}
}

// TestParseShareRefuses sends the lines a garbling backhaul might deliver in place of PEL's share.
func TestParseShareRefuses(t *testing.T) {
	tests := []struct{ name, text string }{
		{"62 digits, a line cut short", pelShare[:62]},
		{"a z among the digits", pelShare[:63] + "z"},
		{"the same share in capitals", strings.ToUpper(pelShare)},
		{"a space at the end", pelShare + " "},
		{"an empty line", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseShare(tt.text); err == nil {
				t.Errorf("parseShare(%q) = a share, nil, want an error", short(tt.text))
			}
		})
	}
}

func TestAgreeHex(t *testing.T) {
	got, err := agreeHex(pelKey(t), ystShare)
	if err != nil || hex.EncodeToString(got) != agreed {
		t.Errorf("agreeHex(PEL's key, %q) = %s, %v, want %s, nil", short(ystShare), short(hex.EncodeToString(got)), err, short(agreed))
	}
}

// TestZeroShare sends 64 zeros: well formed, so parseShare accepts them, but they force an
// all-zero secret, which ECDH refuses.
func TestZeroShare(t *testing.T) {
	zeros := strings.Repeat("0", 64)
	if _, err := parseShare(zeros); err != nil {
		t.Errorf("parseShare(%q) = %v, want a share: the digits are well formed", short(zeros), err)
	}
	got, err := agreeHex(pelKey(t), zeros)
	if err == nil {
		t.Errorf("agreeHex(PEL's key, %q) = %d bytes, nil, want an error", short(zeros), len(got))
	}
	if bytes.Equal(got, make([]byte, 32)) {
		t.Errorf("agreeHex(PEL's key, %q) = 32 zero bytes, want no secret at all", short(zeros))
	}
}
