package main

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// The tests' keys are Cirrine-1's and Pipit-2's seal keys, derived from the test master
// k-2028-07 as lesson 5.4 derives them. The nonce is the fixed one of the test vectors: real
// nonces are always random.
const (
	cirrineSealHex = "427a8ae78a295915acb18a87bdbb40be20d49d27c15dac44e320d4510e070171"
	pipitSealHex   = "53fdd2ccce4f186b0bdeb36626b6e8591ac1c5f59a00a80956c348dbdbaa2c3a"
	nonceHex       = "a1b2c3d4e5f60718293a4b5c"
	command        = "SET_MODE|SAFE"
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

func TestRoundTrip(t *testing.T) {
	key, nonce := fromHex(t, cirrineSealHex), fromHex(t, nonceHex)
	got, err := open(key, nonce, seal(key, nonce, []byte(command)))
	if err != nil || string(got) != command {
		t.Errorf("open(seal(%q)) = %q, %v, want %q, nil", command, got, err, command)
	}
}

func TestTagLength(t *testing.T) {
	key, nonce := fromHex(t, cirrineSealHex), fromHex(t, nonceHex)
	if got := len(seal(key, nonce, []byte(command))); got != len(command)+16 {
		t.Errorf("len(seal(%q)) = %d, want %d: the command's 13 bytes and a 16-byte tag", command, got, len(command)+16)
	}
}

// TestBitFlip plays Mirela, who has no key. She knows that SAFE is bytes 9 to 12 of the
// command, and XORs those locked bytes with SAFE XOR LEOP.
func TestBitFlip(t *testing.T) {
	key, nonce := fromHex(t, cirrineSealHex), fromHex(t, nonceHex)
	flipped := seal(key, nonce, []byte(command))
	if len(flipped) < len(command) {
		t.Fatalf("len(seal(%q)) = %d, want at least %d", command, len(flipped), len(command))
	}
	for i := range 4 {
		flipped[9+i] ^= "SAFE"[i] ^ "LEOP"[i]
	}
	if got, err := open(key, nonce, flipped); err == nil {
		t.Errorf("open(Mirela's flipped bytes) = %q, nil, want an error", got)
	}
}

func TestRefusals(t *testing.T) {
	key, nonce := fromHex(t, cirrineSealHex), fromHex(t, nonceHex)
	sealed := seal(key, nonce, []byte(command))
	if len(sealed) == 0 {
		t.Fatalf("seal(%q) = nothing, want the locked bytes", command)
	}
	changed := bytes.Clone(sealed)
	changed[len(changed)-1] ^= 0x01
	tests := []struct {
		name        string
		key, sealed []byte
	}{
		{"its last byte changed", key, changed},
		{"its last byte cut off", key, sealed[:len(sealed)-1]},
		{"Pipit-2's key", fromHex(t, pipitSealHex), sealed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := open(tt.key, nonce, tt.sealed); err == nil {
				t.Errorf("open(the sealed command, %s) = %q, nil, want an error", tt.name, got)
			}
		})
	}
}
