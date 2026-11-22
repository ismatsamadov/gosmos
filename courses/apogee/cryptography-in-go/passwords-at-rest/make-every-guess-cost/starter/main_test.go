package main

import (
	"encoding/hex"
	"testing"
)

const (
	mirelaPassword = "fails closed, always"
	mirelaSalt     = "mirela test salt" // the tests' readable stand-in for 16 random bytes
	// mirelaKey is PBKDF2-HMAC-SHA256 of mirelaPassword and mirelaSalt at 600,000 iterations,
	// worked out once with crypto/pbkdf2.
	mirelaKey = "6d91ceab779ff574b23820b133a740fa76d8b6a0deb5d6f0da2f94a6508dd938"
)

// short returns a hex value's first 8 digits and an ellipsis when it is longer than 50.
func short(h string) string {
	if len(h) > 50 {
		return h[:8] + "…"
	}
	return h
}

func TestIterations(t *testing.T) {
	if iterations != 600000 {
		t.Errorf("iterations = %d, want 600000", iterations)
	}
}

func TestDerive(t *testing.T) {
	key, err := derive(mirelaPassword, []byte(mirelaSalt))
	if got := hex.EncodeToString(key); err != nil || got != mirelaKey {
		t.Errorf("derive(%q, %q) = %s, %v, want %s, nil", mirelaPassword, mirelaSalt, short(got), err, short(mirelaKey))
	}
}

// TestAnotherSalt derives Mirela's password with a different salt: the key must change too.
func TestAnotherSalt(t *testing.T) {
	const other = "other test salt!"
	key, err := derive(mirelaPassword, []byte(other))
	got := hex.EncodeToString(key)
	if err != nil || len(key) != 32 || got == mirelaKey {
		t.Errorf("derive(%q, %q) = %s (%d bytes), %v, want 32 bytes other than Mirela's key, nil", mirelaPassword, other, short(got), len(key), err)
	}
}
