package main

import (
	"bytes"
	"strings"
	"testing"
)

func passes(t *testing.T, args ...string) (code int, out, errOut string) {
	t.Helper()
	var o, e bytes.Buffer
	code = run(args, &o, &e)
	return code, o.String(), e.String()
}

// The pass docs/apogee/story.md quotes for Season 7's take-over: the model must still give it
// to the second and the tenth of a degree.
func TestThePassTheBibleQuotes(t *testing.T) {
	code, out, errOut := passes(t, "--sat", "furrow-a", "--station", "SQT", "--from", "2029-06-28")
	want := "furrow-a/SQT/24534 2029-06-28T14:02:10Z 2029-06-28T14:12:04Z 47.9\n"
	if code != 0 || !strings.Contains(out, want) {
		t.Fatalf("passes = %d, %q, %q; want a line %q", code, out, errOut, want)
	}
}

// The capstone's investors' week, as docs/apogee/story.md counts it: 933 passes of the five
// satellites fitted in Season 3 over the five original stations, and 213 of Pipit-2's.
func TestTheInvestorsWeek(t *testing.T) {
	week := []string{"--from", "2029-07-23", "--to", "2029-07-30", "--station", "YST,PEL,SQT,LAN,CWR"}
	for _, c := range []struct {
		sats string
		want int
	}{
		{"pipit-1,cirrine-1,brine-3,furrow-a,furrow-b", 933},
		{"pipit-2", 213},
	} {
		code, out, errOut := passes(t, append(week, "--sat", c.sats)...)
		if got := strings.Count(out, "\n"); code != 0 || got != c.want {
			t.Errorf("passes of %s in the investors' week = %d lines (exit %d, %q), want %d", c.sats, got, code, errOut, c.want)
		}
	}
}

func TestPassesComeByAOS(t *testing.T) {
	_, out, _ := passes(t, "--from", "2029-07-23T06:00:00Z", "--to", "2029-07-23T12:00:00Z")
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 100 {
		t.Fatalf("six hours of the whole network = %d passes, want at least 100", len(lines))
	}
	for i := 1; i < len(lines); i++ {
		if strings.Fields(lines[i])[1] < strings.Fields(lines[i-1])[1] {
			t.Fatalf("line %d %q comes after %q, which has a later AOS", i+1, lines[i], lines[i-1])
		}
	}
}

func TestUsageMistakes(t *testing.T) {
	for _, c := range []struct {
		args []string
		want string
	}{
		{nil, "--from is required"},
		{[]string{"--from", "23 July"}, "neither a day"},
		{[]string{"--from", "2029-07-23", "--to", "2029-07-22"}, "--to must be after --from"},
		{[]string{"--from", "2029-07-23", "--sat", "pipit-9"}, `no satellite "pipit-9"; there are pipit-1, pipit-2`},
		{[]string{"--from", "2029-07-23", "--station", "XYZ"}, `no station "XYZ"`},
		{[]string{"--from", "2029-07-23", "extra"}, `unexpected argument "extra"`},
	} {
		code, _, errOut := passes(t, c.args...)
		if code != 2 || !strings.Contains(errOut, c.want) {
			t.Errorf("passes %q = %d, %q; want 2 and %q", c.args, code, errOut, c.want)
		}
	}
}

// Run back from 2029, the fitted satellites fly as Season 3's model had them: its reference
// passes in docs/apogee/story.md's Season 3 section, and Pipit-2's gctl-night pass of Season 2
// (YST, AOS 04:31, 47.2°). The catalog there keeps the two Furrow passes a few seconds shorter
// than the model gives (14:11:52 and 14:58:58), and Season 3's lessons print the catalog's.
func TestSeasonThreesModel(t *testing.T) {
	for _, c := range []struct{ sat, station, day, want string }{
		{"pipit-1", "PEL", "2027-12-06", "pipit-1/PEL/08154 2027-12-06T04:12:05Z 2027-12-06T04:21:27Z 61.4"},
		{"cirrine-1", "YST", "2027-12-06", "cirrine-1/YST/11364 2027-12-06T07:30:00Z 2027-12-06T07:40:04Z 44.2"},
		{"brine-3", "LAN", "2027-12-06", "brine-3/LAN/09213 2027-12-06T09:14:20Z 2027-12-06T09:22:02Z 38.5"},
		{"furrow-a", "SQT", "2027-12-06", "furrow-a/SQT/15864 2027-12-06T14:02:10Z 2027-12-06T14:12:04Z 47.9"},
		{"furrow-b", "SQT", "2027-12-06", "furrow-b/SQT/15862 2027-12-06T14:49:30Z 2027-12-06T14:59:23Z 44.0"},
		{"furrow-a", "SQT", "2028-02-01", "furrow-a/SQT/16731 2028-02-01T14:02:10Z 2028-02-01T14:12:04Z 47.9"},
		{"pipit-2", "YST", "2027-10-30", "pipit-2/YST/13287 2027-10-30T04:31:00Z 2027-10-30T04:40:33Z 47.2"},
	} {
		code, out, errOut := passes(t, "--sat", c.sat, "--station", c.station, "--from", c.day)
		if code != 0 || errOut != "" || !strings.Contains(out, c.want+"\n") {
			t.Errorf("passes of %s over %s on %s = %d, %q\n%s\nwant a line %q", c.sat, c.station, c.day, code, errOut, out, c.want)
		}
	}
}

func TestBeforeTheirFirstPassesSaysSo(t *testing.T) {
	for _, c := range []struct{ sat, from, note string }{
		{"pipit-3", "2028-06-05", "pipit-3 has no passes before 2029-01-01T00:00:00Z"}, // elements from 2029's tracking
		{"pipit-1", "2026-06-05", "pipit-1 has no passes before 2027-01-01T00:00:00Z"},
		{"thistle-2", "2029-03-01", "thistle-2 has no passes before 2029-03-08T10:20:00Z"}, // deployed then
	} {
		code, _, errOut := passes(t, "--sat", c.sat, "--from", c.from)
		if code != 0 || !strings.Contains(errOut, c.note) {
			t.Errorf("passes of %s from %s = %d, %q; want 0 and %q", c.sat, c.from, code, errOut, c.note)
		}
	}
	if _, out, _ := passes(t, "--sat", "pipit-3", "--from", "2028-06-05"); out != "" {
		t.Errorf("passes of pipit-3 in 2028 = %q, want none", out)
	}
}
