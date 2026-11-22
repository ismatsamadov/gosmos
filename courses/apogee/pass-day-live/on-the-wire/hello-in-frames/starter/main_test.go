package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"
	"testing"
	"time"
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

// testEmpty checks that a file of the rehearsal's folder is there and empty.
func testEmpty(t *testing.T, res *simResult, name string) {
	t.Helper()
	lines, ok := res.fileLines(name)
	switch {
	case !ok:
		t.Errorf("%s is not there, want it there and empty", name)
	case len(lines) > 0:
		t.Errorf("%s has %d lines, the first %q; want it empty", name, len(lines), lines[0])
	}
}

// testBooking returns a booking of an eight-minute pass.
func testBooking(id, sat, station string, h, m, s int) Booking {
	aos := simJuly(26, h, m, s)
	return Booking{ID: id, SatID: sat, Station: station, AOS: aos, LOS: aos.Add(8 * 60e9)}
}

func TestBookings(t *testing.T) {
	t.Run("the week's export", func(t *testing.T) {
		res := simWeekRun(t)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testEmpty(t, res, "reports/unplanned.txt")
		want := fmt.Sprintf("bookings: %d, unplanned 0", len(res.c.bookings))
		if !raceOn && want != "bookings: 1493, unplanned 0" {
			t.Fatalf("the week's export has %d bookings, want 1493", len(res.c.bookings))
		}
		testLine(t, res, 2, want)
	})
	t.Run("one of each", func(t *testing.T) {
		c := testThursday("YST", "PEL", "CWR", "QND", "PVT")
		c.bookings = []Booking{
			testBooking("bk-8743", "furrow-h", "YST", 11, 5, 54),
			testBooking("bk-9901", "thistle", "PEL", 10, 15, 0),
			testBooking("bk-9902", "furrow-a", "CWR", 4, 58, 30),
			testBooking("bk-9903", "pipit-1", "YST", 1, 9, 32),
			testBooking("bk-9904", "cirrine-1", "PVT", 0, 17, 43),
			testBooking("bk-9905", "furrow-a", "QND", 0, 2, 58),
			testBooking("bk-9906", "furrow-a", "SQX", 14, 2, 10),
			testBooking("bk-9907", "furrow-h", "YST", 11, 5, 55),
			testBooking("bk-9908", "furrow-a", "SQT", 14, 2, 10),
			{ID: "bk-9909", SatID: "thistle-2", Station: "PEL", AOS: simJuly(27, 1, 2, 3), LOS: simJuly(27, 1, 9, 3)},
		}
		res := simRun(t, c)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		got, _ := res.fileLines("reports/unplanned.txt")
		want := []string{
			"bk-9905 Thu 00:02:58 furrow-a/QND no dish free",
			"bk-9904 Thu 00:17:43 cirrine-1/PVT too low",
			"bk-9903 Thu 01:09:32 pipit-1/YST not servable",
			"bk-9902 Thu 04:58:30 furrow-a/CWR no such pass",
			"bk-9901 Thu 10:15:00 thistle/PEL unknown satellite",
			"bk-9907 Thu 11:05:55 furrow-h/YST no such pass",
			"bk-9906 Thu 14:02:10 furrow-a/SQX unknown station",
		}
		for i := range max(len(got), len(want)) {
			switch {
			case i >= len(got):
				t.Errorf("reports/unplanned.txt ends after %d lines, want %d; the next is %q", len(got), len(want), want[i])
			case i >= len(want):
				t.Errorf("reports/unplanned.txt line %d = %q, want the file to end after %d lines", i+1, got[i], len(want))
			case got[i] != want[i]:
				t.Errorf("reports/unplanned.txt line %d = %q, want %q", i+1, got[i], want[i])
			default:
				continue
			}
			break
		}
		testLine(t, res, 2, "bookings: 8, unplanned 7")
	})
	t.Run("no export", func(t *testing.T) {
		res := simRun(t, testThursday("PEL"))
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testEmpty(t, res, "reports/unplanned.txt")
		testLine(t, res, 2, "bookings: 0, unplanned 0")
	})
	t.Run("a broken export", func(t *testing.T) {
		c := testThursday("PEL")
		c.bookingsText = `[{"id":"bk-8792","sat":"thistle-2"}]` + "\n"
		testRefused(t, simRun(t, c), `bookings.json: json: unknown field "sat"`)
	})
}

// testAt writes an instant for a message: "Thu 16:03:46", "Thu 16:03:46.35" between seconds,
// or with its date when it falls on another day than the instant it is compared with.
func testAt(at, other time.Time) string {
	if at.Format("2006-01-02") == other.Format("2006-01-02") {
		return at.Format("Mon 15:04:05.999")
	}
	return at.Format("Mon 2 Jan 2006 15:04:05.999")
}

// testPassLines is one pass's AOS and LOS lines on the console.
type testPassLines struct {
	aos, los []simEvent
}

// testByPass returns the console's AOS and LOS lines by pass ID.
func testByPass(res *simResult) map[string]*testPassLines {
	out := map[string]*testPassLines{}
	for _, e := range res.events() {
		pl := out[e.PassID]
		if pl == nil {
			pl = &testPassLines{}
			out[e.PassID] = pl
		}
		if e.Kind == "AOS" {
			pl.aos = append(pl.aos, e)
		} else {
			pl.los = append(pl.los, e)
		}
	}
	return out
}

// testOne checks a pass's AOS or LOS lines: exactly one, written at the instant at and
// saying want. It reports the first problem and returns false.
func testOne(t *testing.T, kind, id string, got []simEvent, at time.Time, rest string) bool {
	t.Helper()
	want := fmt.Sprintf("%s %s %s %s", at.Format("Mon 15:04:05"), kind, id, rest)
	switch {
	case len(got) == 0:
		t.Errorf("console has no %s line for %s, want %q", kind, id, want)
	case len(got) > 1:
		t.Errorf("console has %d %s lines for %s, want one: %q", len(got), kind, id, want)
	case !got[0].At.Equal(at):
		t.Errorf("%s line for %s at %s, want %s", kind, id, testAt(got[0].At, at), testAt(at, got[0].At))
	default:
		e := got[0]
		if line := fmt.Sprintf("%s %s %s %s", e.Stamp, e.Kind, e.PassID, e.Rest); line != want {
			t.Errorf("%s line = %q, want %q", kind, line, want)
			return false
		}
		return true
	}
	return false
}

// testFlew checks that every contact of the run's plan at the given stations, all of them
// when none is given, has one AOS line at its AOS and one LOS line at its LOS ending in its
// outcome, and that no other pass has a line. It reports the first problem only.
func testFlew(t *testing.T, res *simResult, outcome func(c Contact) string, stations ...string) {
	t.Helper()
	plan := res.plan()
	byPass := testByPass(res)
	for _, c := range plan.contacts {
		if len(stations) > 0 && !slices.Contains(stations, c.Station) {
			continue
		}
		pl := byPass[c.ID]
		if pl == nil {
			pl = &testPassLines{}
		}
		if !testOne(t, "AOS", c.ID, pl.aos, c.AOS, fmt.Sprintf("dish %d", c.Dish)) ||
			!testOne(t, "LOS", c.ID, pl.los, c.LOS, outcome(c)) {
			return
		}
	}
	for _, e := range res.events() {
		if _, ok := plan.byID[e.PassID]; !ok {
			t.Errorf("console has an %s line for %q, a pass the plan doesn't fly", e.Kind, e.PassID)
			return
		}
	}
}

// testLastLOS returns the latest LOS of the run's plan.
func testLastLOS(res *simResult) time.Time {
	var last time.Time
	for _, c := range res.plan().contacts {
		if c.LOS.After(last) {
			last = c.LOS
		}
	}
	return last
}

// testEnd checks that the console's last line is want, written at the last LOS of the run's
// plan, and that rehearse returned then.
func testEnd(t *testing.T, res *simResult, want string) {
	t.Helper()
	last := testLastLOS(res)
	if n := len(res.console); n == 0 || res.lines()[n-1] != want {
		got := ""
		if n > 0 {
			got = res.lines()[n-1]
		}
		t.Errorf("last console line = %q, want %q", got, want)
	} else if at := res.console[n-1].At; !at.Equal(last) {
		t.Errorf("%q written at %s, want %s, the last LOS", want, testAt(at, last), testAt(last, at))
	}
	if !res.ended.Equal(last) {
		t.Errorf("rehearse() returned at %s, want %s, the last LOS", testAt(res.ended, last), testAt(last, res.ended))
	}
}

// testFlown is every contact's outcome when nothing goes wrong.
func testFlown(Contact) string { return "flown" }

// testWeather returns each contact's outcome in a case's weather: a pass still flying when
// its station's dishes stow is cut short, failed stowed, and one whose AOS comes while they
// are stowed is skipped; every other pass is flown.
func testWeather(c simCase) func(Contact) string {
	return func(ct Contact) string {
		for _, s := range c.stows {
			switch {
			case s.Station != ct.Station:
			case ct.AOS.Before(s.From) && ct.LOS.After(s.From):
				return "failed stowed"
			case !ct.AOS.Before(s.From) && ct.AOS.Before(s.To):
				return "skipped stowed"
			}
		}
		return "flown"
	}
}

// testWeek is each contact's outcome in the investors' week as Ingrid scripts it: on
// Thursday the wind stows Ystmark's dishes from 10:50 to 12:00.
func testWeek(ct Contact) string {
	var storm simCase
	simStorm(&storm)
	return testWeather(storm)(ct)
}

// testDone returns the line that ends a run whose contacts end as outcome says.
func testDone(res *simResult, outcome func(Contact) string) string {
	n := map[string]int{}
	for _, c := range res.plan().contacts {
		kind, _, _ := strings.Cut(outcome(c), " ")
		n[kind]++
	}
	return fmt.Sprintf("week done: %d flown, %d failed, %d skipped", n["flown"], n["failed"], n["skipped"])
}

func TestPasses(t *testing.T) {
	t.Run("the investors' week", func(t *testing.T) {
		res := simWeekRun(t)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testFlew(t, res, testWeek)
	})
	t.Run("two dishes at Ystmark", func(t *testing.T) {
		res := simRun(t, testThursday("YST"))
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		// Dish 2 takes Furrow-M at 10:46:20, while dish 1 still flies Furrow-N.
		aos := simJuly(26, 10, 46, 20)
		m := testByPass(res)["furrow-m/YST/00436"]
		switch {
		case m == nil || len(m.aos) == 0:
			t.Errorf("console has no AOS line for furrow-m/YST/00436, want one at Thu 10:46:20, while dish 1 still flies furrow-n/YST/00436")
			return
		case !m.aos[0].At.Equal(aos):
			t.Errorf("AOS line for furrow-m/YST/00436 at %s, want Thu 10:46:20, while dish 1 still flies furrow-n/YST/00436", testAt(m.aos[0].At, aos))
			return
		}
		testFlew(t, res, testFlown)
	})
	t.Run("the week's end", func(t *testing.T) {
		res := simWeekRun(t)
		want := testDone(res, testWeek)
		if !raceOn && want != "week done: 14283 flown, 2 failed, 12 skipped" {
			t.Fatalf("the reference plan ends %q, want %q", want, "week done: 14283 flown, 2 failed, 12 skipped")
		}
		testEnd(t, res, want)
	})
}

// testNoAnswer is the outcome of a pass whose station never answered: no banner, a wrong
// one, or no WELCOME.
func testNoAnswer(Contact) string { return "failed no_answer" }

// testCalls checks the calls the rehearsal made to a station, or to every station when
// station is "": one per contact, to its station, dialed at its AOS and hung up at its
// LOS, or by its LOS when atLOS is false or the station refused the pass or hung up first.
// It reports the first problem only.
func testCalls(t *testing.T, res *simResult, station string, atLOS bool) {
	t.Helper()
	type key struct {
		station string
		dialed  int64
	}
	calls := map[key][]*simCall{}
	made := 0
	for _, cl := range res.world.calls {
		if station == "" || cl.station == station {
			k := key{cl.station, cl.dialed.UnixNano()}
			calls[k] = append(calls[k], cl)
			made++
		}
	}
	want := 0
	for _, c := range res.plan().contacts {
		if station != "" && c.Station != station {
			continue
		}
		want++
		k := key{c.Station, c.AOS.UnixNano()}
		got := calls[k]
		if len(got) == 0 {
			for _, cl := range res.world.calls {
				if cl.dialed.Equal(c.AOS) && cl.station != c.Station {
					t.Errorf("call for %s at its AOS went to %q, want %q", c.ID, cl.station, c.Station)
					return
				}
			}
			t.Errorf("no call to %s at %s for %s, want one at its AOS", c.Station, c.AOS.Format("Mon 15:04:05"), c.ID)
			return
		}
		i := slices.IndexFunc(got, func(cl *simCall) bool { return cl.hungUp.Equal(c.LOS) })
		if i < 0 {
			i = 0
		}
		cl := got[i]
		calls[k] = slices.Delete(got, i, i+1)
		exact := atLOS && cl.nak == 0 && (cl.closed.IsZero() || !cl.closed.Before(c.LOS))
		switch {
		case cl.hungUp.IsZero():
			t.Errorf("call for %s never hung up, want it hung up at its LOS, %s", c.ID, c.LOS.Format("Mon 15:04:05"))
		case exact && !cl.hungUp.Equal(c.LOS):
			t.Errorf("call for %s hung up at %s, want %s, its LOS", c.ID, testAt(cl.hungUp, c.LOS), testAt(c.LOS, cl.hungUp))
		case !exact && cl.hungUp.After(c.LOS):
			t.Errorf("call for %s hung up at %s, want by %s, its LOS", c.ID, testAt(cl.hungUp, c.LOS), testAt(c.LOS, cl.hungUp))
		default:
			continue
		}
		return
	}
	if made != want {
		t.Errorf("rehearsal made %d calls, want %d: one per contact", made, want)
	}
}

func TestCalls(t *testing.T) {
	t.Run("the investors' week", func(t *testing.T) {
		res := simWeekRun(t)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testCalls(t, res, "", true)
	})
	c := testThursday("PEL", "TIR", "BRM", "CRP")
	c.split = map[string]int{"PEL": 1}
	c.banner = map[string]string{"TIR": "TIR READY PROTO 1\n", "BRM": "BRM READY PROTO 2"}
	c.dialFail = []simSpan{{Station: "CRP"}}
	res := simRun(t, c)
	if res.err != nil {
		t.Fatalf("rehearse() = %v, want nil", res.err)
	}
	t.Run("Port Ellery, a byte at a time", func(t *testing.T) {
		testFlew(t, res, testFlown, "PEL")
		testCalls(t, res, "PEL", true)
	})
	t.Run("Tirrin Rock, LF alone", func(t *testing.T) {
		testFlew(t, res, testFlown, "TIR")
		testCalls(t, res, "TIR", true)
	})
	t.Run("Ilha Brumosa, PROTO 2", func(t *testing.T) {
		testFlew(t, res, testNoAnswer, "BRM")
		testCalls(t, res, "BRM", false)
	})
	t.Run("Corith Point, no line", func(t *testing.T) {
		testFlew(t, res, testNoAnswer, "CRP")
	})
}

func TestSilence(t *testing.T) {
	c := simCase{days: 7, stations: []string{"CAL", "PEL"}, silent: []string{"CAL"}}
	res := simRun(t, c)
	if res.err != nil {
		t.Fatalf("rehearse() = %v, want nil", res.err)
	}
	t.Run("Calidra Atoll, silent", func(t *testing.T) {
		testFlew(t, res, testNoAnswer, "CAL")
		testCalls(t, res, "CAL", true)
	})
	t.Run("Port Ellery beside it", func(t *testing.T) {
		testFlew(t, res, testFlown, "PEL")
		testCalls(t, res, "PEL", true)
		testEnd(t, res, "week done: 614 flown, 500 failed, 0 skipped")
	})
}

// testStopped checks a run stopped at the instant stop. A pass that started before it has
// its AOS line, and its LOS line at its LOS, or at the stop saying stopped when it was still
// flying; a later pass has no line. No call was dialed after the stop, and every call was
// hung up by it, at once when its pass was still flying. The console ends with the stop's
// line, written as rehearse returned. It reports the first problem of each kind.
func testStopped(t *testing.T, res *simResult, stop time.Time) {
	t.Helper()
	byPass := testByPass(res)
	los := map[string]time.Time{} // each contact's LOS, by station and AOS
	for _, c := range res.plan().contacts {
		los[c.Station+" "+c.AOS.String()] = c.LOS
		pl := byPass[c.ID]
		if pl == nil {
			pl = &testPassLines{}
		}
		if !c.AOS.Before(stop) {
			if lines := append(pl.aos, pl.los...); len(lines) > 0 {
				t.Errorf("%s line for %s at %s, want none: its AOS comes after the stop at %s", lines[0].Kind, c.ID, testAt(lines[0].At, stop), testAt(stop, lines[0].At))
				break
			}
			continue
		}
		at, outcome := c.LOS, "flown"
		if c.LOS.After(stop) {
			at, outcome = stop, "stopped"
		}
		if !testOne(t, "AOS", c.ID, pl.aos, c.AOS, fmt.Sprintf("dish %d", c.Dish)) || !testOne(t, "LOS", c.ID, pl.los, at, outcome) {
			break
		}
	}
	for _, cl := range res.world.calls {
		if cl.dialed.After(stop) {
			t.Errorf("call to %s dialed at %s, want none after the stop at %s", cl.station, testAt(cl.dialed, stop), testAt(stop, cl.dialed))
			break
		}
		if cl.late {
			t.Errorf("call to %s dialed at %s with the rehearsal's context ended, want none after the stop", cl.station, testAt(cl.dialed, stop))
			break
		}
		if cl.failed {
			continue
		}
		want := stop
		if end, ok := los[cl.station+" "+cl.dialed.String()]; ok && end.Before(stop) {
			want = end
		}
		if !cl.hungUp.Equal(want) {
			t.Errorf("call to %s dialed at %s hung up at %s, want %s", cl.station, cl.dialed.Format("Mon 15:04:05"), testAt(cl.hungUp, want), testAt(want, cl.hungUp))
			break
		}
	}
	want := "rehearsal stopped at " + stop.Format("Mon 15:04:05")
	lines := res.lines()
	switch n := len(lines); {
	case n == 0 || lines[n-1] != want:
		got := ""
		if n > 0 {
			got = lines[n-1]
		}
		t.Errorf("last console line = %q, want %q", got, want)
	case !res.console[n-1].At.Equal(stop) || !res.ended.Equal(stop):
		t.Errorf("rehearse() returned at %s, want %s, at the stop", testAt(res.ended, stop), testAt(stop, res.ended))
	}
}

func TestStop(t *testing.T) {
	t.Run("Tuesday noon, with a cause", func(t *testing.T) {
		cause := errors.New("Ingrid spotted a wrong setting")
		stop := simJuly(24, 12, 0, 0)
		res := simRun(t, simCase{days: 2, stations: []string{"YST", "TIR"}, stopAt: stop, cause: cause})
		if !errors.Is(res.err, cause) {
			t.Errorf("rehearse() = %v, want the cause %q", res.err, cause)
		}
		testStopped(t, res, stop)
	})
	t.Run("before the first AOS", func(t *testing.T) {
		stop := simJuly(23, 0, 0, 1)
		res := simRun(t, simCase{days: 1, stations: []string{"YST"}, stopAt: stop})
		if !errors.Is(res.err, context.Canceled) {
			t.Errorf("rehearse() = %v, want context.Canceled", res.err)
		}
		testStopped(t, res, stop)
	})
}

// testFrame checks the start of what a station received against want, the frame that
// should come first, part by part: its length, type, seq, body and CRC. It reports the first
// part that differs.
func testFrame(t *testing.T, name string, got, want []byte) {
	t.Helper()
	if len(got) == 0 {
		t.Errorf("%s: the station received nothing, want %d bytes", name, len(want))
		return
	}
	parts := []struct {
		name     string
		from, to int
		why      string
	}{
		{"length", 0, 2, fmt.Sprintf(": %d bytes follow, big-endian", len(want)-2)},
		{"type", 2, 3, ""},
		{"seq", 3, 7, ""},
		{"body", 7, len(want) - 4, ""},
		{"CRC", len(want) - 4, len(want), ", the CRC-32 of the type, seq and body"},
	}
	for _, p := range parts {
		if len(got) < p.to {
			t.Errorf("%s ends after %d bytes, want %d", name, len(got), len(want))
			return
		}
		g, w := got[p.from:p.to], want[p.from:p.to]
		switch {
		case bytes.Equal(g, w):
			continue
		case p.name == "body":
			t.Errorf("%s: body = %q, want %q", name, g, w)
		default:
			t.Errorf("%s: %s = % x, want % x%s", name, p.name, g, w, p.why)
		}
		return
	}
}

// testHellos checks that every call the rehearsal made to a station at a contact's AOS said
// HELLO for that contact and its dish, or for the other contact with that AOS when the
// station's two dishes both start one then. It reports the first problem only.
func testHellos(t *testing.T, res *simResult, station string) {
	t.Helper()
	for _, cl := range res.world.calls {
		if cl.station != station || cl.failed {
			continue
		}
		var want []string
		for _, c := range res.plan().contacts {
			if c.Station == station && c.AOS.Equal(cl.dialed) {
				want = append(want, fmt.Sprintf("%s %d", c.ID, c.Dish))
			}
		}
		switch {
		case len(want) == 0 || slices.Contains(want, cl.hello):
			continue
		case cl.hello == "":
			id, _, _ := strings.Cut(want[0], " ")
			t.Errorf("call for %s said no HELLO, want HELLO %q", id, want[0])
		default:
			id, _, _ := strings.Cut(want[0], " ")
			t.Errorf("call for %s said HELLO %q, want %q", id, cl.hello, want[0])
		}
		return
	}
}

func TestHello(t *testing.T) {
	c := testThursday("PEL", "TIR", "YST")
	c.glue = []string{"TIR"}
	c.stows = []Stow{{Station: "YST", From: simJuly(26, 10, 50, 0), To: simJuly(26, 12, 0, 0), Reason: "wind lock"}}
	res := simRun(t, c)
	if res.err != nil {
		t.Fatalf("rehearse() = %v, want nil", res.err)
	}
	t.Run("Port Ellery's HELLO", func(t *testing.T) {
		want, _ := hex.DecodeString(strings.ReplaceAll("00 1e 01 00 00 00 00 74 68 69 73 74 6c 65 2d 32 "+
			"2f 50 45 4c 2f 30 32 31 33 36 20 31 ba 2c bf 54", " ", ""))
		aos := simJuly(26, 16, 3, 46)
		i := slices.IndexFunc(res.world.calls, func(cl *simCall) bool { return cl.station == "PEL" && cl.dialed.Equal(aos) })
		if i < 0 {
			t.Errorf("no call to PEL at Thu 16:03:46 for thistle-2/PEL/02136, want one at its AOS")
		} else {
			testFrame(t, "HELLO for thistle-2/PEL/02136", res.world.calls[i].received, want)
		}
		testHellos(t, res, "PEL")
		testFlew(t, res, testFlown, "PEL")
	})
	t.Run("Tirrin Rock, banner and WELCOME in one write", func(t *testing.T) {
		testHellos(t, res, "TIR")
		testFlew(t, res, testFlown, "TIR")
	})
	t.Run("Ystmark's wind lock", func(t *testing.T) {
		testFlew(t, res, testWeather(c), "YST")
	})
	t.Run("Cape Wren, no WELCOME", func(t *testing.T) {
		c := testThursday("CWR")
		c.bannerOnly = true
		res := simRun(t, c)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testFlew(t, res, testNoAnswer)
		testCalls(t, res, "CWR", true)
	})
	t.Run("a pass left out of the archive", func(t *testing.T) {
		// Without pipit-2/TIR/22864 (08:15:11 to 08:23:23, dish 1), the rehearsal gives
		// furrow-b/TIR/24954 (08:17:25 to 08:25:39) dish 1, where the station flies it on 2.
		c := testThursday("TIR")
		var archive strings.Builder
		for _, line := range strings.SplitAfter(string(simArchiveText(simWindow(c.from, c.days), c.stations)), "\n") {
			if !strings.HasPrefix(line, "pipit-2/TIR/22864 ") {
				archive.WriteString(line)
			}
		}
		c.passes = archive.String()
		res := simRun(t, c)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		byPass := testByPass(res)
		for _, ct := range res.plan().contacts {
			pl := byPass[ct.ID]
			if pl == nil {
				pl = &testPassLines{}
			}
			dish, outcome := ct.Dish, "flown"
			switch ct.ID {
			case "pipit-2/TIR/22864":
				if lines := append(pl.aos, pl.los...); len(lines) > 0 {
					t.Errorf("console has an %s line for pipit-2/TIR/22864, want none: the archive doesn't have it", lines[0].Kind)
					return
				}
				continue
			case "furrow-b/TIR/24954":
				dish, outcome = 1, "failed refused"
			}
			if !testOne(t, "AOS", ct.ID, pl.aos, ct.AOS, fmt.Sprintf("dish %d", dish)) || !testOne(t, "LOS", ct.ID, pl.los, ct.LOS, outcome) {
				return
			}
		}
	})
}
