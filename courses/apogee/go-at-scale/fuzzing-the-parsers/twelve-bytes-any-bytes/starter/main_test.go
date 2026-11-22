package main

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// beacon returns the bytes a beacon's hex digits spell.
func beacon(digits string) []byte {
	b, _ := hex.DecodeString(digits)
	return b
}

// Beacons heard on Monday 16 April: Pipit-1's over Port Ellery, Pipit-2's over Tirrin Rock.
var (
	pipit1 = beacon("415001de3b001efa00290038")
	pipit2 = beacon("415002a0f7001f540070007f")
)

// safeDecode calls decodeBeacon on b and fails the test as decodeBeacon(...) panicked: <value>, instead of ending the test binary.
func safeDecode(t *testing.T, b []byte) (Beacon, error) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("decodeBeacon(% x) panicked: %v", b, r)
		}
	}()
	return decodeBeacon(b)
}

// FuzzBeacon checks decodeBeacon on every seed and every file under testdata/fuzz/FuzzBeacon: it
// never panics, and a beacon it accepts encodes back to exactly the bytes it came from.
func FuzzBeacon(f *testing.F) {
	f.Add(pipit1)
	f.Add(pipit2)
	f.Add(append(bytes.Clone(pipit1), 0x00))
	f.Add(beacon("415101de3b001efa00290039"))
	f.Fuzz(func(t *testing.T, b []byte) {
		bc, err := safeDecode(t, b)
		if err != nil {
			return
		}
		if got := encodeBeacon(bc); !bytes.Equal(got, b) {
			t.Errorf("encodeBeacon(decodeBeacon(% x)) = % x, want the same bytes back", b, got)
		}
	})
}

func TestDecodeBeacon(t *testing.T) {
	t.Run("Pipit-1's beacon over Port Ellery", func(t *testing.T) {
		bc, err := safeDecode(t, pipit1)
		if err != nil {
			t.Fatalf("decodeBeacon(Pipit-1's beacon) error = %v, want nil", err)
		}
		if want := (Beacon{Sat: 1, Seq: 56891, Mode: 0, BattMV: 7930, TempDeciC: 41}); bc != want {
			t.Errorf("decodeBeacon(Pipit-1's beacon) = %+v, want %+v", bc, want)
		}
	})
	tests := []struct {
		name string
		b    []byte
		want string
	}{
		{"Pipit-2's beacon, cut at Tirrin Rock's LOS after 7 bytes", pipit2[:7], "beacon: want 12 bytes, got 7"},
		{"no bytes at all", nil, "beacon: want 12 bytes, got 0"},
		{"a beacon and one byte more", append(bytes.Clone(pipit1), 0x00), "beacon: want 12 bytes, got 13"},
		{"sync bytes AQ", beacon("415101de3b001efa00290039"), "beacon: no sync"},
		{"a check byte one off", beacon("415001de3b001efa00290039"), "beacon: bad check byte"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := safeDecode(t, tt.b); err == nil || err.Error() != tt.want {
				t.Errorf("decodeBeacon(% x) error = %v, want %q", tt.b, err, tt.want)
			}
		})
	}
}
