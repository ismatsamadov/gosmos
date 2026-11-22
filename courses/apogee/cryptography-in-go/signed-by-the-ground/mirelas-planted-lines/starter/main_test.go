package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"maps"
	"slices"
	"testing"
)

// The lab bench key signs Mirela's lines; flight keys never appear in code.
var (
	benchKey  = []byte("lab bench key, never flown")
	forgerKey = []byte("a key Mirela made up herself")
)

// sigOf returns the signature key gives a canonical command, as a sig= field carries it.
func sigOf(key []byte, canonical string) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(canonical))
	return hex.EncodeToString(mac.Sum(nil))
}

// errAny stands for a refusal of any kind: a line that doesn't parse or a field canonical refuses.
var errAny = errors.New("any refusal")

func TestPlantedLines(t *testing.T) {
	sig119 := sigOf(benchKey, "pipit-2|119|SET_MODE|NOMINAL")
	sig120 := sigOf(benchKey, "pipit-2|120|SET_MODE|SAFE")
	line1 := "CMD pipit-2 118 SET_MODE SAFE sig=512545bcb43a5bb6ccb9f9c8511f4cc8c5fedeb38c226376af5b5aaae60e3b4d"
	lines := []struct {
		name, line string
		want       error // nil when the line is accepted
	}{
		{"1, good", line1, nil},
		{"2, good", "CMD pipit-2 119 SET_MODE NOMINAL sig=" + sig119, nil},
		{"3, an argument changed", "CMD pipit-2 120 SET_MODE LEOP sig=" + sig120, ErrForged},
		{"4, the number changed", "CMD pipit-2 121 SET_MODE SAFE sig=" + sig120, ErrForged},
		{"5, a bar inside an argument", "CMD pipit-2 120 SET_MODE SAFE|REBOOT sig=" + sigOf(benchKey, "pipit-2|120|SET_MODE|SAFE|REBOOT"), errAny},
		{"6, no sig= field", "CMD pipit-2 120 SET_MODE SAFE", errAny},
		{"7, 63 hex digits", "CMD pipit-2 120 SET_MODE SAFE sig=" + sig120[:63], errAny},
		{"8, signed with another key", "CMD pipit-2 120 SET_MODE SAFE sig=" + sigOf(forgerKey, "pipit-2|120|SET_MODE|SAFE"), ErrForged},
		{"9, sig= not last", "CMD pipit-2 120 sig=" + sig120 + " SET_MODE SAFE", errAny},
		{"10, the largest number, forged", "CMD pipit-2 18446744073709551615 SET_MODE SAFE sig=" + sigOf(forgerKey, "pipit-2|18446744073709551615|SET_MODE|SAFE"), ErrForged},
		{"11, line 1 again", line1, ErrReplay},
		{"12, an empty line", "", errAny},
	}
	var g Guard
	for _, l := range lines {
		before := maps.Clone(g.last)
		c, err := checkLine(l.line, benchKey, &g)
		if l.want == nil {
			if err != nil || c.SatID != "pipit-2" {
				t.Errorf("checkLine(line %s) = %+v, %v, want its command, nil", l.name, c, err)
			}
			continue
		}
		if err == nil || (l.want != errAny && !errors.Is(err, l.want)) {
			t.Errorf("checkLine(line %s) error = %v, want %v", l.name, err, l.want)
		}
		if !maps.Equal(before, g.last) {
			t.Errorf("checkLine(line %s) moved the guard from %v to %v, want it unchanged", l.name, before, g.last)
		}
	}
	if want := map[string]uint64{"pipit-2": 119}; !maps.Equal(g.last, want) {
		t.Errorf("the guard after all twelve lines = %v, want %v", g.last, want)
	}
}

func TestAcceptedCommand(t *testing.T) {
	var g Guard
	line := "CMD pipit-2 119 SET_MODE NOMINAL sig=" + sigOf(benchKey, "pipit-2|119|SET_MODE|NOMINAL")
	c, err := checkLine(line, benchKey, &g)
	if err != nil || c.SatID != "pipit-2" || c.Seq != 119 || c.Op != "SET_MODE" || !slices.Equal(c.Args, []string{"NOMINAL"}) || c.Sig != line[len(line)-64:] {
		t.Errorf("checkLine(%q) = %+v, %v, want pipit-2's command 119 SET_MODE NOMINAL with its signature, nil", line[:33]+"…", c, err)
	}
}
