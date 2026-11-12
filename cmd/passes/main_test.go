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

func TestBeforeTheElementsSaysSo(t *testing.T) {
	code, out, errOut := passes(t, "--sat", "pipit-1", "--from", "2028-06-05")
	if code != 0 || out != "" || !strings.Contains(errOut, "pipit-1's elements start at 2029-01-01T00:00:00Z") {
		t.Errorf("passes before the epoch = %d, %q, %q; want 0, nothing and a note", code, out, errOut)
	}
}
