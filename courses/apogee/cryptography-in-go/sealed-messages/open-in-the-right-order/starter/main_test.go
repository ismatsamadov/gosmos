package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

// The tests' keys are Cirrine-1's two keys and Pipit-2's seal key, all from the test master
// k-2028-07 (lesson 5.4). pingLine is the sealed PING of the test vectors, made with the fixed
// test nonce, which real seals never use.
const (
	cirrineSignHex = "18668bfe29024937ff84a61df34e8a0310a4702a74834b969f56e2c6ccd00213"
	cirrineSealHex = "427a8ae78a295915acb18a87bdbb40be20d49d27c15dac44e320d4510e070171"
	pipitSealHex   = "53fdd2ccce4f186b0bdeb36626b6e8591ac1c5f59a00a80956c348dbdbaa2c3a"
	pingLine       = "SEAL 2 cirrine-1 4471 k-2028-07 a1b2c3d4e5f60718293a4b5c c1f2a9798d38081d764bbcfb5e35e2598d2584a4"
)

// fromHex decodes one of the tests' hex constants.
func fromHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("hex.DecodeString(%q): %v", s, err)
	}
	return b
}

// cirrine returns the tests' Cirrine-1 at 07:30:05 on 25 August, its last accepted number last.
func cirrine(t *testing.T, last uint64) *Satellite {
	t.Helper()
	return &Satellite{
		ID:   "cirrine-1",
		Keys: map[string]SatKeys{"k-2028-07": {Sign: fromHex(t, cirrineSignHex), Seal: fromHex(t, cirrineSealHex)}},
		Last: last,
		Now:  func() time.Time { return time.Date(2028, 8, 25, 7, 30, 5, 0, time.UTC) },
	}
}

// sealFor seals plaintext under a seal key, as uplinkd does, with a random nonce.
func sealFor(t *testing.T, key []byte, sat string, seq uint64, keyID, plaintext string) Sealed {
	t.Helper()
	aead, err := newAEAD(key)
	if err != nil {
		t.Fatal(err)
	}
	s := Sealed{Version: 2, SatID: sat, Seq: seq, KeyID: keyID, Nonce: make([]byte, 12)}
	rand.Read(s.Nonce)
	s.Ciphertext = aead.Seal(nil, s.Nonce, []byte(plaintext), s.AD())
	return s
}

func TestFirstPing(t *testing.T) {
	sat := cirrine(t, 4470)
	c, err := sat.Accept(pingLine)
	if err != nil || c.SatID != "cirrine-1" || c.Seq != 4471 || c.Op != "PING" || len(c.Args) != 0 {
		t.Errorf("Accept(the sealed PING) = %s %d %s %q, %v, want cirrine-1 4471 PING [], nil", c.SatID, c.Seq, c.Op, c.Args, err)
	}
	if sat.Last != 4471 {
		t.Errorf("Last after the sealed PING = %d, want 4471", sat.Last)
	}
}

// TestRefusals gives Cirrine-1, which has just accepted the PING numbered 4471, one bad line
// each. Every refusal must give its reason and leave Last at 4471.
func TestRefusals(t *testing.T) {
	cirrineSeal := fromHex(t, cirrineSealHex)
	changed := sealFor(t, cirrineSeal, "cirrine-1", 4472, "k-2028-07", "SET_MODE|SAFE")
	changed.Ciphertext[0] ^= 0x01
	tests := []struct {
		name     string
		line     string
		want     error
		wantName string
	}{
		{"the same PING again", pingLine, ErrReplay, "ErrReplay"},
		{"a line for pipit-2", formatSealed(sealFor(t, fromHex(t, pipitSealHex), "pipit-2", 4472, "k-2028-07", "PING")), ErrWrongSat, "ErrWrongSat"},
		{"a line under k-2027-03", formatSealed(sealFor(t, cirrineSeal, "cirrine-1", 4472, "k-2027-03", "PING")), ErrUnknownKey, "ErrUnknownKey"},
		{"one changed ciphertext byte", formatSealed(changed), ErrForged, "ErrForged"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sat := cirrine(t, 4471)
			if _, err := sat.Accept(tt.line); !errors.Is(err, tt.want) {
				t.Errorf("Accept(%s) = %v, want an error wrapping %s", tt.name, err, tt.wantName)
			}
			if sat.Last != 4471 {
				t.Errorf("Last after %s = %d, want 4471: a refusal records nothing", tt.name, sat.Last)
			}
		})
	}
}

// TestKeyState gives Cirrine-1 a good line under k-2028-07 once the satellite has revoked that
// key, and once the key has retired at 07:30:00, five seconds before the stand-in's clock.
func TestKeyState(t *testing.T) {
	line := formatSealed(sealFor(t, fromHex(t, cirrineSealHex), "cirrine-1", 4472, "k-2028-07", "PING"))
	tests := []struct {
		name     string
		change   func(k *SatKeys)
		want     error
		wantName string
	}{
		{"revoked", func(k *SatKeys) { k.Revoked = true }, ErrRevoked, "ErrRevoked"},
		{"retired at 07:30:00", func(k *SatKeys) { k.Retire = time.Date(2028, 8, 25, 7, 30, 0, 0, time.UTC) }, ErrRetired, "ErrRetired"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sat := cirrine(t, 4471)
			keys := sat.Keys["k-2028-07"]
			tt.change(&keys)
			sat.Keys["k-2028-07"] = keys
			if _, err := sat.Accept(line); !errors.Is(err, tt.want) {
				t.Errorf("Accept(a line under k-2028-07, %s) = %v, want an error wrapping %s", tt.name, err, tt.wantName)
			}
			if sat.Last != 4471 {
				t.Errorf("Last after a line under k-2028-07, %s = %d, want 4471: a refusal records nothing", tt.name, sat.Last)
			}
		})
	}
}

// TestForgedMaximum plays Mirela: she renumbers the sealed PING to the largest number there is,
// then uplinkd sends the next real command.
func TestForgedMaximum(t *testing.T) {
	sat := cirrine(t, 4471)
	forged := strings.Replace(pingLine, " 4471 ", " "+strconv.FormatUint(math.MaxUint64, 10)+" ", 1)
	if _, err := sat.Accept(forged); !errors.Is(err, ErrForged) {
		t.Errorf("Accept(the PING renumbered 18446744073709551615) = %v, want an error wrapping ErrForged", err)
	}
	if sat.Last != 4471 {
		t.Errorf("Last after the forged line = %d, want 4471: a refusal records nothing", sat.Last)
	}
	next := formatSealed(sealFor(t, fromHex(t, cirrineSealHex), "cirrine-1", 4472, "k-2028-07", "SET_MODE|SAFE"))
	if c, err := sat.Accept(next); err != nil || c.Seq != 4472 {
		t.Errorf("Accept(the real command numbered 4472, after the forgery) = %d, %v, want 4472, nil", c.Seq, err)
	}
}
