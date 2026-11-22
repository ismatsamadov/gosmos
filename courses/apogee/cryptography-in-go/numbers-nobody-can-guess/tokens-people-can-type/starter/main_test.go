package main

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

// base32 is the alphabet of rand.Text: A to Z, then 2 to 7.
const base32 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

func TestTokenShape(t *testing.T) {
	tok := newSessionToken()
	if len(tok) < 26 {
		t.Errorf("len(newSessionToken()) = %d, want at least 26", len(tok))
	}
	for i, r := range tok {
		if !strings.ContainsRune(base32, r) {
			t.Errorf("newSessionToken() = %q, with %q at %d, want only A to Z and 2 to 7", tok, r, i)
			break
		}
	}
}

func TestTokensDiffer(t *testing.T) {
	a, b := newSessionToken(), newSessionToken()
	if a == b {
		t.Errorf("newSessionToken() = %q twice in a row, want two different tokens", a)
	}
}

// TestClockReader plays someone who knows, give or take a second, when you logged in, and
// reads your token as a time.
func TestClockReader(t *testing.T) {
	from := time.Now().Add(-time.Second).UnixNano()
	tok := newSessionToken()
	to := time.Now().Add(time.Second).UnixNano()
	if ns, err := strconv.ParseInt(tok, 10, 64); err == nil && ns >= from && ns <= to {
		t.Errorf("newSessionToken() = %q, the login time %s, want a token no clock gives away", tok, time.Unix(0, ns).UTC().Format(time.RFC3339Nano))
	}
}
