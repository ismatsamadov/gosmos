package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// The tests' stand-ins for k-2028-07 and k-2028-10: 32 readable bytes each.
var (
	july    = Key{ID: "k-2028-07", Secret: []byte("uplink test master key, 32 bytes")}
	october = Key{ID: "k-2028-10", Secret: []byte("uplink next master key, 32 bytes")}
)

// verbs are the ways a debug line can print a value.
var verbs = []string{"%v", "%+v", "%s", "%#v"}

// short cuts a printed value longer than 50 characters to its first 40 and an ellipsis.
func short(s string) string {
	if len(s) > 50 {
		return s[:40] + "…"
	}
	return s
}

// leak says how out shows one of the test keys' secrets: as text, as decimal bytes, as hex or
// as Go hex. It returns "" when out shows none of them.
func leak(out string) string {
	for _, k := range []Key{july, october} {
		s := k.Secret
		forms := []struct{ how, text string }{
			{"as text", string(s)},
			{"as decimal bytes", fmt.Sprintf("%d %d %d %d", s[0], s[1], s[2], s[3])},
			{"as hex", fmt.Sprintf("%x", s[:4])},
			{"as Go hex", fmt.Sprintf("%#x, %#x, %#x, %#x", s[0], s[1], s[2], s[3])},
		}
		for _, f := range forms {
			if strings.Contains(out, f.text) {
				return f.how
			}
		}
	}
	return ""
}

// checkPrints prints v, described as what, with every verb, and fails when a secret shows or
// one of the IDs doesn't.
func checkPrints(t *testing.T, what string, v any, ids ...string) {
	t.Helper()
	for _, verb := range verbs {
		out := fmt.Sprintf(verb, v)
		if how := leak(out); how != "" {
			t.Errorf("fmt.Sprintf(%q, %s) = %q, which shows a secret %s, want no secret", verb, what, short(out), how)
		}
		for _, id := range ids {
			if !strings.Contains(out, id) {
				t.Errorf("fmt.Sprintf(%q, %s) = %q, want it to show the ID %s", verb, what, short(out), id)
			}
		}
	}
}

func TestKey(t *testing.T) {
	checkPrints(t, "key", july, "k-2028-07")
	if got, want := fmt.Sprint(july), "k-2028-07 (secret hidden)"; got != want {
		t.Errorf("fmt.Sprint(key) = %q, want %q", short(got), want)
	}
	if got, want := fmt.Sprintf("%#v", july), `Key{ID: "k-2028-07", Secret: hidden}`; got != want {
		t.Errorf("fmt.Sprintf(\"%%#v\", key) = %q, want %q", short(got), want)
	}
}

func TestKeySlice(t *testing.T) {
	checkPrints(t, "[]Key{july, october}", []Key{july, october}, "k-2028-07", "k-2028-10")
}

// logLine is how a debug line might carry a key: in an exported field of a struct.
type logLine struct {
	Sat string
	Key Key
}

func TestStruct(t *testing.T) {
	checkPrints(t, "a struct holding a key", logLine{Sat: "cirrine-1", Key: july}, "k-2028-07")
}

// TestRing prints a *KeyRing. Its versions sit in an unexported field, where fmt can't call
// Key's methods, so the ring needs String and GoString methods of its own.
func TestRing(t *testing.T) {
	ring := &KeyRing{versions: []KeyVersion{
		{Key: july, Activate: time.Date(2028, 8, 11, 0, 0, 0, 0, time.UTC)},
		{Key: october, Activate: time.Date(2028, 10, 17, 0, 0, 0, 0, time.UTC)},
	}}
	checkPrints(t, "ring", ring, "k-2028-07", "k-2028-10")
	if got, want := fmt.Sprint(ring), "key ring: k-2028-07, k-2028-10"; got != want {
		t.Errorf("fmt.Sprint(ring) = %q, want %q", short(got), want)
	}
}
