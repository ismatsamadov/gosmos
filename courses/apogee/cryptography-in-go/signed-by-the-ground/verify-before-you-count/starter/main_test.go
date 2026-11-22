package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"
	"strings"
	"testing"
)

// The lab bench key signs the tests' commands; flight keys never appear in code.
var benchKey = []byte("lab bench key, never flown")

// setMode returns Pipit-2's SET_MODE command numbered seq, not yet signed.
func setMode(seq uint64, mode string) Command {
	return Command{SatID: "pipit-2", Seq: seq, Op: "SET_MODE", Args: []string{mode}}
}

// signed returns c with the signature that key gives its canonical form.
func signed(t *testing.T, c Command, key []byte) Command {
	t.Helper()
	text, err := canonical(c)
	if err != nil {
		t.Fatalf("canonical(%v) error = %v, want nil", c, err)
	}
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(text))
	c.Sig = hex.EncodeToString(mac.Sum(nil))
	return c
}

func TestForgedMaximum(t *testing.T) {
	var g Guard
	forged := setMode(math.MaxUint64, "SAFE")
	forged.Sig = strings.Repeat("0", 64)
	if err := g.Accept(forged, benchKey); !errors.Is(err, ErrForged) {
		t.Errorf("Accept(forged command %d) = %v, want ErrForged", forged.Seq, err)
	}
	if err := g.Accept(signed(t, setMode(119, "NOMINAL"), benchKey), benchKey); err != nil {
		t.Errorf("Accept(signed command 119) after the forgery = %v, want nil", err)
	}
}

func TestReplay(t *testing.T) {
	var g Guard
	c := signed(t, setMode(118, "SAFE"), benchKey)
	if err := g.Accept(c, benchKey); err != nil {
		t.Fatalf("Accept(command 118) = %v, want nil", err)
	}
	const want = "replayed sequence number: 118, last 118"
	if err := g.Accept(c, benchKey); !errors.Is(err, ErrReplay) || err.Error() != want {
		t.Errorf("Accept(command 118) again = %v, want %q", err, want)
	}
}

func TestInOrder(t *testing.T) {
	var g Guard
	for _, seq := range []uint64{119, 120} {
		if err := g.Accept(signed(t, setMode(seq, "NOMINAL"), benchKey), benchKey); err != nil {
			t.Errorf("Accept(command %d) = %v, want nil", seq, err)
		}
	}
}

func TestRefusalChangesNothing(t *testing.T) {
	var g Guard
	if err := g.Accept(signed(t, setMode(118, "SAFE"), benchKey), benchKey); err != nil {
		t.Fatalf("Accept(command 118) = %v, want nil", err)
	}
	forged := signed(t, setMode(500, "LEOP"), []byte("someone else's key"))
	if err := g.Accept(forged, benchKey); !errors.Is(err, ErrForged) {
		t.Errorf("Accept(command 500 signed with another key) = %v, want ErrForged", err)
	}
	if last := g.last["pipit-2"]; last != 118 {
		t.Errorf("Pipit-2's last number after the refusal = %d, want 118", last)
	}
}
