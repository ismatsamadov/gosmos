package main

import (
	"encoding/hex"
	"testing"
)

// testMaster is the tests' stand-in for k-2028-07: 32 readable bytes.
var testMaster = Key{ID: "k-2028-07", Secret: []byte("uplink test master key, 32 bytes")}

// short returns a hex value's first 8 digits and an ellipsis when it is longer than 50.
func short(h string) string {
	if len(h) > 50 {
		return h[:8] + "…"
	}
	return h
}

// TestStolenKey plays Mirela at the lab bench. She has read Pipit-2's key out of the spare
// flight computer, and signs a command for Cirrine-1 with it.
func TestStolenKey(t *testing.T) {
	stolen, err := signKey(testMaster, "pipit-2")
	if err != nil {
		t.Fatalf("signKey(testMaster, %q) = %v, want a key", "pipit-2", err)
	}
	cirrine, err := signKey(testMaster, "cirrine-1")
	if err != nil {
		t.Fatalf("signKey(testMaster, %q) = %v, want a key", "cirrine-1", err)
	}
	const cmd = "cirrine-1|4402|SET_MODE|SAFE"
	if verify(cmd, hex.EncodeToString(tag(stolen, cmd)), cirrine) {
		t.Errorf("verify(%q, its signature by Pipit-2's key, Cirrine-1's key) = true, want false: a key read out of Pipit-2 must not command Cirrine-1", cmd)
	}
}

func TestSignKey(t *testing.T) {
	tests := []struct{ sat, want string }{
		{"pipit-2", "82f945638cbdd9a034dfab2bc578fd7a720d73e5701d3abcf561cd825847da7c"},
		{"cirrine-1", "18668bfe29024937ff84a61df34e8a0310a4702a74834b969f56e2c6ccd00213"},
	}
	for _, tt := range tests {
		key, err := signKey(testMaster, tt.sat)
		if got := hex.EncodeToString(key); err != nil || got != tt.want {
			t.Errorf("signKey(testMaster, %q) = %s (%d bytes), %v, want %s (32 bytes), nil", tt.sat, short(got), len(key), err, short(tt.want))
		}
	}
}

func TestEmptySatellite(t *testing.T) {
	if key, err := signKey(testMaster, ""); err == nil {
		t.Errorf("signKey(testMaster, \"\") = %d bytes, nil, want an error", len(key))
	}
}
