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

func TestSatKey(t *testing.T) {
	tests := []struct{ purpose, sat, want string }{
		{"sign", "cirrine-1", "18668bfe29024937ff84a61df34e8a0310a4702a74834b969f56e2c6ccd00213"},
		{"seal", "cirrine-1", "427a8ae78a295915acb18a87bdbb40be20d49d27c15dac44e320d4510e070171"},
		{"seal", "pipit-2", "53fdd2ccce4f186b0bdeb36626b6e8591ac1c5f59a00a80956c348dbdbaa2c3a"},
	}
	for _, tt := range tests {
		key, err := satKey(testMaster, tt.purpose, tt.sat)
		if got := hex.EncodeToString(key); err != nil || got != tt.want {
			t.Errorf("satKey(testMaster, %q, %q) = %s (%d bytes), %v, want %s (32 bytes), nil", tt.purpose, tt.sat, short(got), len(key), err, short(tt.want))
		}
	}
}

func TestSatKeyRefuses(t *testing.T) {
	tests := []struct{ name, purpose, sat string }{
		{"a third purpose", "wrap", "cirrine-1"},
		{"Thistle, which re-entered in 2027", "sign", "thistle"},
		{"no satellite at all", "sign", ""},
		{"a display name, not an ID", "seal", "Cirrine-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if key, err := satKey(testMaster, tt.purpose, tt.sat); err == nil {
				t.Errorf("satKey(testMaster, %q, %q) = %d bytes, nil, want an error", tt.purpose, tt.sat, len(key))
			}
		})
	}
}
