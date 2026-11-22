package main

import (
	"errors"
	"io/fs"
	"testing"
)

// testThursday is a case flying Thursday 26 July at the given stations.
func testThursday(stations ...string) simCase {
	return simCase{from: simJuly(26, 0, 0, 0), days: 1, stations: stations}
}

// testFirst checks that the console's first Write is exactly want and a newline.
func testFirst(t *testing.T, res *simResult, want string) {
	t.Helper()
	if res.err != nil {
		t.Fatalf("rehearse() = %v, want nil", res.err)
	}
	if len(res.console) == 0 {
		t.Fatalf("console is empty, want its first line %q", want)
	}
	if got := res.console[0].Text; got != want+"\n" {
		t.Errorf("first console write = %q, want %q", got, want+"\n")
	}
}

// testRefused checks that rehearse returned an error with the text want and printed
// nothing.
func testRefused(t *testing.T, res *simResult, want string) {
	t.Helper()
	if res.err == nil || res.err.Error() != want {
		t.Errorf("rehearse() = %v, want %q", res.err, want)
	}
	if len(res.console) > 0 {
		t.Errorf("console = %q, want nothing after an error", res.lines())
	}
}

func TestSettings(t *testing.T) {
	t.Run("only the week", func(t *testing.T) {
		res := simWeekRun(t)
		want := "rehearsal: 2029-07-23, 7 days, 20 stations"
		if raceOn {
			want = "rehearsal: 2029-07-23, 7 days, 4 stations" // the four of simRaceStations
		}
		testFirst(t, res, want)
	})
	accepted := []struct {
		name string
		c    simCase
		text string
		want string
	}{
		{"one day, two stations", testThursday("TIR", "YST"), `{"week":"2029-07-26","days":1,"stations":["TIR","YST"]}`,
			"rehearsal: 2029-07-26, 1 day, 2 stations"},
		{"seven days, one station", simCase{from: simWeek, days: 7, stations: []string{"PEL"}}, `{"week":"2029-07-23","days":7,"stations":["PEL"]}`,
			"rehearsal: 2029-07-23, 7 days, 1 station"},
		{"stations in any order", testThursday("YST", "PEL", "VRL"), `{"stations":["VRL","PEL","YST"],"days":1,"week":"2029-07-26"}`,
			"rehearsal: 2029-07-26, 1 day, 3 stations"},
	}
	for _, tc := range accepted {
		t.Run(tc.name, func(t *testing.T) {
			tc.c.settings = tc.text + "\n"
			testFirst(t, simRun(t, tc.c), tc.want)
		})
	}
	refused := []struct {
		name, text, want string
	}{
		{"an unknown field", `{"week":"2029-07-26","dayz":1}`, `rehearsal.json: json: unknown field "dayz"`},
		{"a second object", `{"week":"2029-07-26"}{"week":"2029-07-27"}`, `rehearsal.json: more data after the settings`},
		{"no week", `{"days":1,"stations":["PEL"]}`, `rehearsal.json: week is required`},
		{"not a date", `{"week":"2029-7-26"}`, `rehearsal.json: week "2029-7-26" is not a date like 2006-01-02`},
		{"a time of day", `{"week":"2029-07-26T00:00:00Z"}`, `rehearsal.json: week "2029-07-26T00:00:00Z" is not a date like 2006-01-02`},
		{"zero days", `{"week":"2029-07-26","days":0}`, `rehearsal.json: days 0 is not from 1 to 7`},
		{"eight days", `{"week":"2029-07-26","days":8}`, `rehearsal.json: days 8 is not from 1 to 7`},
		{"an unknown station", `{"week":"2029-07-26","stations":["YST","YSX"]}`, `rehearsal.json: unknown station "YSX"`},
	}
	for _, tc := range refused {
		t.Run(tc.name, func(t *testing.T) {
			c := testThursday("PEL")
			c.settings = tc.text + "\n"
			testRefused(t, simRun(t, c), tc.want)
		})
	}
	t.Run("no settings", func(t *testing.T) {
		c := testThursday("PEL")
		c.settings = "-"
		res := simRun(t, c)
		if !errors.Is(res.err, fs.ErrNotExist) {
			t.Errorf("rehearse() = %v, want an error that wraps fs.ErrNotExist", res.err)
		}
		if len(res.console) > 0 {
			t.Errorf("console = %q, want nothing after an error", res.lines())
		}
	})
}
