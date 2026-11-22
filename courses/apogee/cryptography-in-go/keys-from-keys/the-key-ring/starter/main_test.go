package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

const (
	julyHex = "75706c696e6b2074657374206d6173746572206b65792c203332206279746573" // "uplink test master key, 32 bytes"
	nextHex = "75706c696e6b206e657874206d6173746572206b65792c203332206279746573" // "uplink next master key, 32 bytes"
	// augustFile is uplink-keys.txt as uplinkd holds it in August, with the tests' master.
	augustFile = "# uplink-keys.txt: uplinkd only, mode 0600\n" +
		"# <id> <activate> <retire or -> <secret> [revoked]\n" +
		"\n" +
		"k-2028-07 2028-08-11T00:00:00Z - " + julyHex + "\n"
)

func TestLoadRing(t *testing.T) {
	ring, err := loadRing(strings.NewReader(augustFile))
	if err != nil || ring == nil {
		t.Fatalf("loadRing(August's file) = %v, %v, want a ring, nil", ring, err)
	}
	vs := ring.Versions()
	if len(vs) != 1 {
		t.Fatalf("loadRing(August's file) holds %d versions, want 1: k-2028-07", len(vs))
	}
	v, activate := vs[0], time.Date(2028, 8, 11, 0, 0, 0, 0, time.UTC)
	if v.ID != "k-2028-07" || !v.Activate.Equal(activate) || !v.Retire.IsZero() || v.Revoked || hex.EncodeToString(v.Secret) != julyHex {
		t.Errorf("loadRing(August's file) holds %v, activating %s, retiring %v, revoked %v, want k-2028-07 with the file's secret, activating 2028-08-11T00:00:00Z, not retiring, not revoked", v.Key, v.Activate.Format(time.RFC3339), v.Retire, v.Revoked)
	}
}

func TestLoadRingRefuses(t *testing.T) {
	version := func(id, activate, secret string) string {
		return id + " " + activate + " - " + secret + "\n"
	}
	tests := []struct {
		name, file string
		line       int
	}{
		{"a duplicate ID on line 3", "# July, twice\n" + version("k-2028-07", "2028-08-11T00:00:00Z", julyHex) + version("k-2028-07", "2028-10-17T00:00:00Z", nextHex), 3},
		{"a 63-digit secret on line 2", "# cut short\n" + version("k-2028-07", "2028-08-11T00:00:00Z", julyHex[:63]), 2},
		{"a 13th month on line 1", version("k-2028-07", "2028-13-01T00:00:00Z", julyHex), 1},
		{"three fields on line 4", "# first note\n\n# second note\nk-2028-07 2028-08-11T00:00:00Z " + julyHex + "\n", 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ring, err := loadRing(strings.NewReader(tt.file))
			if want := fmt.Sprintf("line %d: ", tt.line); err == nil || !strings.HasPrefix(err.Error(), want) {
				t.Errorf("loadRing(a file with %s) = %v, %v, want an error starting %q", tt.name, ring, err, want)
			}
		})
	}
}

func TestSatKey(t *testing.T) {
	secret, _ := hex.DecodeString(julyHex)
	ring := &KeyRing{}
	if err := ring.add(KeyVersion{Key: Key{ID: "k-2028-07", Secret: secret}, Activate: time.Date(2028, 8, 11, 0, 0, 0, 0, time.UTC)}); err != nil {
		t.Fatalf("ring.add(k-2028-07) = %v, want nil", err)
	}
	const sealCirrine = "427a8ae78a295915acb18a87bdbb40be20d49d27c15dac44e320d4510e070171"
	key, err := ring.SatKey("k-2028-07", "seal", "cirrine-1")
	if got := hex.EncodeToString(key); err != nil || got != sealCirrine {
		t.Errorf("ring.SatKey(%q, %q, %q) = %d bytes, %v, want 427a8ae7… (32 bytes), nil", "k-2028-07", "seal", "cirrine-1", len(key), err)
	}
	key, err = ring.SatKey("k-2027-03", "sign", "pipit-2")
	if !errors.Is(err, ErrUnknownKey) || err.Error() != `unknown key: "k-2027-03"` || key != nil {
		t.Errorf("ring.SatKey(%q, %q, %q) = %d bytes, %v, want nil, ErrUnknownKey (unknown key: \"k-2027-03\")", "k-2027-03", "sign", "pipit-2", len(key), err)
	}
	if key, err := ring.SatKey("k-2028-07", "wrap", "cirrine-1"); err == nil {
		t.Errorf("ring.SatKey(%q, %q, %q) = %d bytes, nil, want satKey's error", "k-2028-07", "wrap", "cirrine-1", len(key))
	}
}
