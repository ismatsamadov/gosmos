package main

import (
	"bytes"
	"testing"
)

// The tests' stand-in for k-2028-07: the 32 ASCII bytes of testMaster, readable on purpose.
const (
	testMaster    = "uplink test master key, 32 bytes"
	testMasterHex = "75706c696e6b2074657374206d6173746572206b65792c203332206279746573"
)

// short returns s, or its first 8 characters and an ellipsis when it is longer than 50.
func short(s string) string {
	if len(s) > 50 {
		return s[:8] + "…"
	}
	return s
}

func TestNewKey(t *testing.T) {
	k := newKey("k-2028-07")
	if k.ID != "k-2028-07" {
		t.Errorf("newKey(%q).ID = %q, want %q", "k-2028-07", k.ID, "k-2028-07")
	}
	if len(k.Secret) != 32 {
		t.Fatalf("len(newKey(%q).Secret) = %d, want 32", "k-2028-07", len(k.Secret))
	}
	if again := newKey("k-2028-07"); bytes.Equal(k.Secret, again.Secret) {
		t.Errorf("newKey(%q).Secret = %x… twice in a row, want two different secrets", "k-2028-07", k.Secret[:4])
	}
}

func TestParseSecret(t *testing.T) {
	got, err := parseSecret(testMasterHex)
	if err != nil || string(got) != testMaster {
		t.Errorf("parseSecret(%q) = %q, %v, want %q, nil", short(testMasterHex), got, err, testMaster)
	}
}

func TestParseSecretRefuses(t *testing.T) {
	tests := []struct{ name, text string }{
		{"62 digits, a line cut short", testMasterHex[:62]},
		{"66 digits, one byte too many", testMasterHex + "21"},
		{"a g among the digits", "g" + testMasterHex[1:]},
		{"an empty line", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := parseSecret(tt.text); err == nil {
				t.Errorf("parseSecret(%q) = %d bytes, nil, want an error", short(tt.text), len(got))
			}
		})
	}
}
