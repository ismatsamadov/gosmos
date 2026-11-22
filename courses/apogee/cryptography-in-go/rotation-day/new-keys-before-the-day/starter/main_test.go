package main

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"
	"time"
)

// at returns the instant of an RFC 3339 time.
func at(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// The two versions in uplinkd's ring since Monday 9 October, with the season's test secrets.
var (
	k07 = KeyVersion{
		Key:      Key{ID: "k-2028-07", Secret: []byte("uplink test master key, 32 bytes")},
		Activate: at("2028-08-11T00:00:00Z"),
		Retire:   at("2028-10-24T00:00:00Z"),
	}
	k10 = KeyVersion{
		Key:      Key{ID: "k-2028-10", Secret: []byte("uplink next master key, 32 bytes")},
		Activate: at("2028-10-17T00:00:00Z"),
	}
)

// cirrine returns Cirrine-1's stand-in on Thursday 12 October, at its 07:30 pass over Ystmark: holding
// its own k-2028-07 keys, its last command numbered 4500.
func cirrine(t *testing.T) *Satellite {
	t.Helper()
	sign, err := satKey(k07.Key, "sign", "cirrine-1")
	if err != nil {
		t.Fatalf("satKey: %v", err)
	}
	seal, err := satKey(k07.Key, "seal", "cirrine-1")
	if err != nil {
		t.Fatalf("satKey: %v", err)
	}
	return &Satellite{
		ID:   "cirrine-1",
		Keys: map[string]SatKeys{"k-2028-07": {Sign: sign, Seal: seal, Retire: k07.Retire}},
		Last: 4500,
		Now:  func() time.Time { return at("2028-10-12T07:30:05Z") },
	}
}

// loadKey calls loadKeyCommand for Cirrine-1 and stops the test unless it returns a line.
func loadKey(t *testing.T, ring *KeyRing, seq uint64, next KeyVersion) string {
	t.Helper()
	line, err := loadKeyCommand(ring, "cirrine-1", seq, next)
	if err != nil || line == "" {
		t.Fatalf("loadKeyCommand(ring, %q, %d, %s) = %q, %v, want a SEAL line", "cirrine-1", seq, next.ID, line, err)
	}
	return line
}

// short shortens hex longer than 50 digits to its first 8 and an ellipsis.
func short(s string) string {
	if len(s) > 50 {
		return s[:8] + "…"
	}
	return s
}

// header returns a SEAL line's readable fields, its hex shortened.
func header(line string) string {
	f := strings.Fields(line)
	if len(f) > 5 {
		f = f[:5]
	}
	return strings.Join(f, " ") + " …"
}

func TestLoadKey(t *testing.T) {
	ring := &KeyRing{versions: []KeyVersion{k07, k10}}
	sat := cirrine(t)
	line := loadKey(t, ring, 4501, k10)
	if c, err := sat.Receive(line); err != nil || c.Op != "LOAD_KEY" {
		t.Fatalf("Cirrine-1: Receive(%q) = %q, %v, want LOAD_KEY, nil", header(line), c.Op, err)
	}
	for _, k := range []struct{ purpose, got, want string }{
		{"sign", hex.EncodeToString(sat.Keys["k-2028-10"].Sign), "3c280cbe585cfec33b8ac6693e6e0802fc0c532efc2af548cb1b3a85054418ba"},
		{"seal", hex.EncodeToString(sat.Keys["k-2028-10"].Seal), "a70fbff3ac2b05cbf0a66c3fc0eb76ceb1ec55269766703f715f4ec43c78663d"},
	} {
		if k.got != k.want {
			t.Errorf("Cirrine-1's k-2028-10 %s key after LOAD_KEY = %s, want %s", k.purpose, short(k.got), short(k.want))
		}
	}
	for _, sealed := range []struct {
		id  string
		seq uint64
	}{{"k-2028-10", 4502}, {"k-2028-07", 4503}} {
		ping, err := sealCommand(ring, sealed.id, Command{SatID: "cirrine-1", Seq: sealed.seq, Op: "PING"})
		if err != nil {
			t.Fatalf("sealCommand(PING with %s) = %v, want a line", sealed.id, err)
		}
		if _, err := sat.Receive(ping); err != nil {
			t.Errorf("Cirrine-1: Receive(PING sealed with %s) = %v, want nil", sealed.id, err)
		}
	}
}

func TestKeysStayHidden(t *testing.T) {
	ring := &KeyRing{versions: []KeyVersion{k07, k10}}
	line := loadKey(t, ring, 4501, k10)
	for _, purpose := range []string{"sign", "seal"} {
		key, err := satKey(k10.Key, purpose, "cirrine-1")
		if err != nil {
			t.Fatalf("satKey: %v", err)
		}
		if strings.Contains(line, hex.EncodeToString(key)) || strings.Contains(line, string(key)) {
			t.Errorf("loadKeyCommand(...) = %q, carrying k-2028-10's %s key readable, want it sealed", header(line), purpose)
		}
	}
	if f := strings.Fields(line); len(f) != 7 || f[4] != "k-2028-07" {
		t.Errorf("loadKeyCommand(...) = %q, want it sealed under k-2028-07, the key Cirrine-1 already holds", header(line))
	}
}

// TestImpostor plays someone who knows the key IDs but not the secrets, and tries to load a key of
// their own into Cirrine-1.
func TestImpostor(t *testing.T) {
	fake := KeyVersion{Key: Key{ID: "k-2028-07", Secret: []byte("another master key, not apogee's")}, Activate: k07.Activate, Retire: k07.Retire}
	theirs := KeyVersion{Key: Key{ID: "k-2028-10", Secret: []byte("another next key, also not ours.")}, Activate: k10.Activate}
	sat := cirrine(t)
	line := loadKey(t, &KeyRing{versions: []KeyVersion{fake, theirs}}, 4501, theirs)
	if _, err := sat.Receive(line); !errors.Is(err, ErrForged) {
		t.Errorf("Cirrine-1: Receive(the impostor's LOAD_KEY) = %v, want ErrForged", err)
	}
	if _, loaded := sat.Keys["k-2028-10"]; loaded || sat.Last != 4500 {
		t.Errorf("after the impostor's LOAD_KEY: k-2028-10 loaded %v, last %d, want false, 4500", loaded, sat.Last)
	}
}

func TestReplay(t *testing.T) {
	ring := &KeyRing{versions: []KeyVersion{k07, k10}}
	sat := cirrine(t)
	if _, err := sat.Receive(loadKey(t, ring, 4501, k10)); err != nil {
		t.Fatalf("Cirrine-1: Receive(the first LOAD_KEY) = %v, want nil", err)
	}
	if _, err := sat.Receive(loadKey(t, ring, 4501, k10)); !errors.Is(err, ErrReplay) {
		t.Errorf("Cirrine-1: Receive(a second LOAD_KEY numbered 4501) = %v, want ErrReplay", err)
	}
}
