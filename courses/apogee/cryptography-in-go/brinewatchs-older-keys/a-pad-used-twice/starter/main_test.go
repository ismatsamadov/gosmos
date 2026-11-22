package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"testing"
)

// The two secrets that went to Brinewatch on the pad (Season 3's test values): its API key in
// November, then its webhook secret in January.
const (
	november = "brinewatch.0f1e2d3c4b5a69788796a5b4c3d2e1f0"
	january  = "brinewatch webhook secret"
)

// testPad returns the tests' copy of Brinewatch's pad: 256 bytes, the SHA-256 of
// "brinewatch pad 0" to "brinewatch pad 7", one after another. A real pad is random bytes;
// this one is readable on purpose.
func testPad() []byte {
	var pad []byte
	for i := range 8 {
		sum := sha256.Sum256(fmt.Appendf(nil, "brinewatch pad %d", i))
		pad = append(pad, sum[:]...)
	}
	return pad
}

// mustSeal seals msg with p, stopping the test if Seal refuses.
func mustSeal(t *testing.T, p *Pad, msg string) []byte {
	t.Helper()
	sealed, err := p.Seal([]byte(msg))
	if err != nil {
		t.Fatalf("Seal(%d bytes) = %v, want nil", len(msg), err)
	}
	if len(sealed) != len(msg) {
		t.Fatalf("Seal(%d bytes) = %d bytes, want %d", len(msg), len(sealed), len(msg))
	}
	return sealed
}

// TestSomeoneKeptBoth plays whoever kept both emails and knows November's secret: XORing the two
// sealed secrets with November's must not give January's.
func TestSomeoneKeptBoth(t *testing.T) {
	p := &Pad{key: testPad()}
	nov := mustSeal(t, p, november)
	jan := mustSeal(t, p, january)
	read := make([]byte, len(january))
	for i := range read {
		read[i] = nov[i] ^ jan[i] ^ november[i]
	}
	if string(read) == january {
		t.Errorf("both sealed secrets XORed with November's = %q, want nothing readable", read)
	}
}

// TestUsedOnce seals three runs of 32 zero bytes, which come out as the very pad bytes each seal
// used: they must be three stretches of the pad, one after another.
func TestUsedOnce(t *testing.T) {
	pad := testPad()
	p := &Pad{key: pad}
	for i := range 3 {
		got := mustSeal(t, p, string(make([]byte, 32)))
		if from := bytes.Index(pad, got); from != 32*i {
			t.Errorf("Seal(32 zero bytes), call %d = the pad's bytes from %d, want from %d", i+1, from, 32*i)
		}
	}
}

// TestBrinewatchOpens opens both secrets with Brinewatch's own copy of the pad, in the order
// they were sealed.
func TestBrinewatchOpens(t *testing.T) {
	apogee, brinewatch := &Pad{key: testPad()}, &Pad{key: testPad()}
	for _, secret := range []string{november, january} {
		sealed := mustSeal(t, apogee, secret)
		got, err := brinewatch.Seal(sealed)
		if err != nil || string(got) != secret {
			t.Errorf("Brinewatch's Seal(sealed %q) = %q, %v, want %q, nil", secret, got, err, secret)
		}
	}
}

// TestRunsOut gives a pad of 50 bytes: November's 43 fit, January's 25 don't, and a refusal
// spends nothing.
func TestRunsOut(t *testing.T) {
	p := &Pad{key: testPad()[:50]}
	mustSeal(t, p, november)
	const want = "pad has 7 bytes left, need 25"
	if _, err := p.Seal([]byte(january)); err == nil || err.Error() != want {
		t.Fatalf("Seal(January's 25 bytes) with 7 left = %v, want the error %q", err, want)
	}
	if _, err := p.Seal([]byte("7 bytes")); err != nil {
		t.Errorf("Seal(7 bytes) after the refusal = %v, want nil: a refusal spends nothing", err)
	}
	if _, err := p.Seal([]byte("!")); err == nil {
		t.Errorf("Seal(1 byte) once all 50 are spent = nil, want %q", "pad has 0 bytes left, need 1")
	}
}
