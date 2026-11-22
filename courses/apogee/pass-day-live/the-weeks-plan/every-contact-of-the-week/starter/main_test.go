package main

import (
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"
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

// testPlanFiles checks plan/<day>.txt for every day of the run's window against the plan
// passplan 2's rule makes of the same archive, and that plan/ holds nothing else.
func testPlanFiles(t *testing.T, res *simResult) {
	t.Helper()
	plan := res.plan()
	days := plan.days()
	for _, day := range days {
		name := "plan/" + day + ".txt"
		got, ok := res.fileLines(name)
		if !ok {
			continue // the check of plan/'s names below reports it
		}
		want := plan.dayLines(day)
		for i := range max(len(got), len(want)) {
			switch {
			case i >= len(got):
				t.Errorf("%s ends after %d lines, want %d; the next is %q", name, len(got), len(want), want[i])
			case i >= len(want):
				t.Errorf("%s line %d = %q, want the file to end after %d lines", name, i+1, got[i], len(want))
			case got[i] != want[i]:
				t.Errorf("%s line %d = %q, want %q", name, i+1, got[i], want[i])
			default:
				continue
			}
			break
		}
	}
	var files []string
	for _, day := range days {
		files = append(files, day+".txt")
	}
	if got := res.names("plan"); !slices.Equal(got, files) {
		t.Errorf("plan/ holds %q, want %q", got, files)
	}
}

// testLine checks the console's line n, counted from 0.
func testLine(t *testing.T, res *simResult, n int, want string) {
	t.Helper()
	lines := res.lines()
	if n >= len(lines) {
		t.Errorf("console has no line %d, want %q", n+1, want)
		return
	}
	if lines[n] != want {
		t.Errorf("console line %d = %q, want %q", n+1, lines[n], want)
	}
}

func TestPlan(t *testing.T) {
	t.Run("the investors' week", func(t *testing.T) {
		res := simWeekRun(t)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testPlanFiles(t, res)
		plan := res.plan()
		want := fmt.Sprintf("plan: %d contacts, %d dropped", len(plan.contacts), plan.dropped)
		if !raceOn && want != "plan: 14297 contacts, 28680 dropped" {
			t.Fatalf("the reference plan says %q, want 14297 contacts and 28680 dropped", want)
		}
		testLine(t, res, 1, want)
	})
	t.Run("past midnight", func(t *testing.T) {
		res := simWeekRun(t)
		for _, c := range res.plan().contacts {
			if c.AOS.Day() == c.LOS.Day() {
				continue
			}
			line := fmt.Sprintf("%s %s %d %s", c.AOS.Format("15:04:05"), c.LOS.Format("15:04:05"), c.Dish, c.ID)
			if lines, _ := res.fileLines("plan/" + c.AOS.Format("2006-01-02") + ".txt"); !slices.Contains(lines, line) {
				t.Errorf("plan/%s.txt has no %q: a contact stays in its AOS's day", c.AOS.Format("2006-01-02"), line)
			}
			break
		}
	})
	t.Run("a day at Port Ellery", func(t *testing.T) {
		res := simRun(t, testThursday("PEL"))
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testPlanFiles(t, res)
		plan := res.plan()
		testLine(t, res, 1, fmt.Sprintf("plan: %d contacts, %d dropped", len(plan.contacts), plan.dropped))
	})
	t.Run("a bad archive line", func(t *testing.T) {
		c := testThursday("PEL")
		c.passes = "thistle-2/PEL/02136 2029-07-26T16:03:46Z 2029-07-26T16:09:59Z 11.7\n" +
			"pipit-2/PEL/22870 2029-07-26T17:33:18Z 2029-07-26T17:43:01Z 61.0\n" +
			"pipit-2/PEL/22871 2029-07-26T19:09:12Z 17.2\n"
		res := simRun(t, c)
		if res.err == nil || !strings.HasPrefix(res.err.Error(), "archive line 3: ") {
			t.Errorf("rehearse() = %v, want an error that starts %q", res.err, "archive line 3: ")
		}
		if len(res.console) > 0 {
			t.Errorf("console = %q, want nothing after an error", res.lines())
		}
	})
	t.Run("no archive", func(t *testing.T) {
		c := testThursday("PEL")
		c.passes = "-"
		res := simRun(t, c)
		if !errors.Is(res.err, fs.ErrNotExist) || len(res.console) > 0 {
			t.Errorf("rehearse() = %v, console %q; want an error that wraps fs.ErrNotExist, and nothing printed", res.err, res.lines())
		}
	})
}
