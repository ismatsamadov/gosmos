package main

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
	"testing/cryptotest"
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
// outcome, as testCommanded says it, and that no other pass has a line. It reports the first
// problem only.
func testFlew(t *testing.T, res *simResult, outcome func(c Contact) string, stations ...string) {
	t.Helper()
	outcome = testCommanded(res, outcome)
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
// are stowed is skipped; a pass flying when its station's line dies, which ends before the
// line is back, has lost it; every other pass is flown.
func testWeather(c simCase) func(Contact) string {
	return func(ct Contact) string {
		for _, q := range c.quiet {
			if q.Station == ct.Station && ct.AOS.Before(q.From) && ct.LOS.After(q.From) && ct.LOS.Before(q.To) {
				return "failed line_lost"
			}
		}
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
// Thursday Ystmark's fiber dies at 03:31, and the wind stows its dishes from 10:50 to 12:00.
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

// testCommanded returns each contact's outcome once commands go up: the outcome given, or
// failed refused when its satellite refused a command from the pass. A copy sent again of a
// command that had already run, refused as replayed, doesn't count: it ran.
func testCommanded(res *simResult, outcome func(Contact) string) func(Contact) string {
	refused := map[string]bool{} // pass IDs
	for _, s := range res.world.sats {
		ran := map[string]bool{} // "<pass ID> <seq>" of each command that ran
		for _, cmd := range s.commands {
			if cmd.Refused == "" {
				ran[fmt.Sprint(cmd.PassID, " ", cmd.Seq)] = true
			}
		}
		for _, cmd := range s.commands {
			again := strings.HasPrefix(cmd.Refused, "replayed sequence number") && ran[fmt.Sprint(cmd.PassID, " ", cmd.Seq)]
			if cmd.Refused != "" && !again {
				refused[cmd.PassID] = true
			}
		}
	}
	return func(ct Contact) string {
		if refused[ct.ID] {
			return "failed refused"
		}
		return outcome(ct)
	}
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
		if want := testDone(res, testWeek); !raceOn && want != "week done: 14282 flown, 3 failed, 12 skipped" {
			t.Fatalf("the reference plan ends %q, want %q", want, "week done: 14282 flown, 3 failed, 12 skipped")
		}
		testEnd(t, res, testDone(res, testCommanded(res, testWeek)))
	})
}

// testNoAnswer is the outcome of a pass whose station never answered: no banner, a wrong
// one, or no WELCOME.
func testNoAnswer(Contact) string { return "failed no_answer" }

// testDialer returns a function that finds the contact of the run's plan a call was made
// for: the one its HELLO named, if that one was in progress at its station then; or else
// the one there whose AOS is the instant it was dialed, or the earliest one in progress.
func testDialer(res *simResult) func(cl *simCall) (Contact, bool) {
	plan := res.plan()
	byStation := map[string][]Contact{} // each station's contacts, by AOS
	for _, c := range plan.contacts {
		byStation[c.Station] = append(byStation[c.Station], c)
	}
	return func(cl *simCall) (Contact, bool) {
		during := func(c Contact) bool {
			return c.Station == cl.station && !cl.dialed.Before(c.AOS) && cl.dialed.Before(c.LOS)
		}
		if c, ok := plan.byID[cl.passID]; ok && during(c) {
			return c, true
		}
		cs := byStation[cl.station]
		i := sort.Search(len(cs), func(i int) bool { return cs[i].AOS.After(cl.dialed) })
		var found []Contact // latest AOS first
		for j := i - 1; j >= 0 && cs[j].AOS.After(cl.dialed.Add(-time.Hour)); j-- {
			if cs[j].AOS.Equal(cl.dialed) {
				return cs[j], true
			}
			if during(cs[j]) {
				found = append(found, cs[j])
			}
		}
		if len(found) == 0 {
			return Contact{}, false
		}
		return found[len(found)-1], true
	}
}

// testCalls checks the calls the rehearsal made to a station, or to every station when
// station is "": each made for a contact in progress there, the first at the contact's AOS,
// every one hung up by its LOS, and a pass's only call at its LOS, unless atLOS is false or
// the station refused the pass or hung up first. It reports the first problem only.
func testCalls(t *testing.T, res *simResult, station string, atLOS bool) {
	t.Helper()
	dialedFor := testDialer(res)
	byPass := map[string][]*simCall{}
	var stray []*simCall // calls made when no contact was in progress at their station
	for _, cl := range res.world.calls {
		if station != "" && cl.station != station {
			continue
		}
		if c, ok := dialedFor(cl); ok {
			byPass[c.ID] = append(byPass[c.ID], cl)
		} else {
			stray = append(stray, cl)
		}
	}
	for _, c := range res.plan().contacts {
		if station != "" && c.Station != station {
			continue
		}
		calls := byPass[c.ID]
		if len(calls) == 0 || !calls[0].dialed.Equal(c.AOS) {
			for _, cl := range res.world.calls {
				if cl.dialed.Equal(c.AOS) && cl.station != c.Station {
					t.Errorf("call for %s at its AOS went to %q, want %q", c.ID, cl.station, c.Station)
					return
				}
			}
			t.Errorf("no call to %s at %s for %s, want one at its AOS", c.Station, c.AOS.Format("Mon 15:04:05"), c.ID)
			return
		}
		for _, cl := range calls {
			switch {
			case cl.failed:
			case cl.hungUp.IsZero():
				t.Errorf("call for %s dialed at %s never hung up, want it hung up by its LOS, %s", c.ID, testAt(cl.dialed, c.LOS), c.LOS.Format("Mon 15:04:05"))
				return
			case cl.hungUp.After(c.LOS):
				t.Errorf("call for %s hung up at %s, want by %s, its LOS", c.ID, testAt(cl.hungUp, c.LOS), testAt(c.LOS, cl.hungUp))
				return
			}
		}
		cl := calls[0]
		exact := len(calls) == 1 && atLOS && cl.nak == 0 && (cl.closed.IsZero() || !cl.closed.Before(c.LOS))
		if exact && !cl.hungUp.Equal(c.LOS) {
			t.Errorf("call for %s hung up at %s, want %s, its LOS", c.ID, testAt(cl.hungUp, c.LOS), testAt(c.LOS, cl.hungUp))
			return
		}
	}
	if len(stray) > 0 {
		cl := stray[0]
		t.Errorf("call to %s dialed at %s, when no contact there was in progress", cl.station, cl.dialed.Format("Mon 15:04:05.999"))
	}
}

// testArchive returns passes.txt for a case's window and stations, except that at the
// station of each pass given it holds that pass alone.
func testArchive(c simCase, ids ...string) string {
	alone := map[string]bool{} // the stations that keep only the passes given
	for _, id := range ids {
		alone[strings.Split(id, "/")[1]] = true
	}
	var b strings.Builder
	for _, line := range strings.SplitAfter(string(simArchiveText(simWindow(c.from, c.days), c.stations)), "\n") {
		id, _, _ := strings.Cut(line, " ")
		if f := strings.Split(id, "/"); len(f) == 3 && alone[f[1]] && !slices.Contains(ids, id) {
			continue
		}
		b.WriteString(line)
	}
	return b.String()
}

// testUnanswered checks a pass whose station never takes it: its AOS and LOS lines, the
// LOS line ending failed no_answer; its first call at its AOS and each later one after the
// wait testRedials wants; every call ended by its LOS. It reports the first problem only.
func testUnanswered(t *testing.T, res *simResult, id string) {
	t.Helper()
	c := res.plan().byID[id]
	pl := testByPass(res)[id]
	if pl == nil {
		pl = &testPassLines{}
	}
	if !testOne(t, "AOS", id, pl.aos, c.AOS, fmt.Sprintf("dish %d", c.Dish)) || !testOne(t, "LOS", id, pl.los, c.LOS, "failed no_answer") {
		return
	}
	dialedFor := testDialer(res)
	var calls []*simCall
	for _, cl := range res.world.calls {
		if d, ok := dialedFor(cl); ok && d.ID == id {
			calls = append(calls, cl)
		}
	}
	if len(calls) == 0 || !calls[0].dialed.Equal(c.AOS) {
		t.Errorf("no call to %s at %s for %s, want one at its AOS", c.Station, c.AOS.Format("Mon 15:04:05"), id)
		return
	}
	for _, cl := range calls {
		if !cl.failed && (cl.hungUp.IsZero() || cl.hungUp.After(c.LOS)) {
			t.Errorf("call for %s dialed at %s hung up at %s, want by %s, its LOS", id, testAt(cl.dialed, c.LOS), testAt(cl.hungUp, c.LOS), testAt(c.LOS, cl.hungUp))
			return
		}
	}
	testRedials(t, res, c)
}

func TestCalls(t *testing.T) {
	t.Run("the investors' week", func(t *testing.T) {
		res := simWeekRun(t)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testCalls(t, res, "", true)
	})
	c := testThursday("PEL", "TIR")
	c.split = map[string]int{"PEL": 1}
	c.banner = map[string]string{"TIR": "TIR READY PROTO 1\n"}
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
	// A call that fails is dialed again after a wait, so Ilha Brumosa and Corith Point
	// fly one pass each, and Ilha Brumosa's run stops 30 s into it.
	t.Run("Ilha Brumosa, PROTO 2", func(t *testing.T) {
		c := testThursday("BRM")
		c.banner = map[string]string{"BRM": "BRM READY PROTO 2"}
		c.passes = testArchive(c, "furrow-j/BRM/00442")
		c.stopAt = simJuly(26, 21, 29, 8)
		res := simRun(t, c)
		for _, cl := range res.world.calls {
			if cl.hello != "" {
				t.Errorf("call to BRM dialed at %s said HELLO after the banner %q, want it hung up", cl.dialed.Format("Mon 15:04:05.999"), "BRM READY PROTO 2")
				return
			}
		}
		testRedials(t, res, res.plan().byID["furrow-j/BRM/00442"])
	})
	t.Run("Corith Point, no line", func(t *testing.T) {
		c := testThursday("CRP")
		c.dialFail = []simSpan{{Station: "CRP"}}
		c.passes = testArchive(c, "brine-20/CRP/01080")
		testUnanswered(t, simRun(t, c), "brine-20/CRP/01080")
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
		testIdle(t, res, "CAL")
	})
	t.Run("Port Ellery beside it", func(t *testing.T) {
		testFlew(t, res, testFlown, "PEL")
		testCalls(t, res, "PEL", true)
		// 614 flown and 500 failed, less the Port Ellery passes whose command was refused.
		testEnd(t, res, testDone(res, testCommanded(res, func(ct Contact) string {
			if ct.Station == "CAL" {
				return "failed no_answer"
			}
			return "flown"
		})))
	})
}

// testStopped checks a run stopped at the instant stop. A pass that started before it has
// its AOS line, and its LOS line at its LOS, flown as testCommanded says it, or at the stop
// saying stopped when it was still flying; a later pass has no line. No call was dialed after the stop, and every call was
// hung up by it, at once when its pass was still flying. The console ends with the stop's
// line, written as rehearse returned. It reports the first problem of each kind.
func testStopped(t *testing.T, res *simResult, stop time.Time) {
	t.Helper()
	byPass := testByPass(res)
	flown := testCommanded(res, testFlown)
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
		at, outcome := c.LOS, flown(c)
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
		testIdle(t, res, "CWR")
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
		byPass, flown := testByPass(res), testCommanded(res, testFlown)
		for _, ct := range res.plan().contacts {
			pl := byPass[ct.ID]
			if pl == nil {
				pl = &testPassLines{}
			}
			dish, outcome := ct.Dish, flown(ct)
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

// testHealth checks reports/health.txt: a line per satellite of the network, in fleet order,
// each the newest good frame that reached the rehearsal or "<id> never heard". It checks the
// satellites given, or all of them when none is, and reports the first problem only.
func testHealth(t *testing.T, res *simResult, sats ...string) {
	t.Helper()
	fleet := simNet().Satellites
	lines, ok := res.fileLines("reports/health.txt")
	if !ok {
		t.Errorf("reports/health.txt is not there, want it there with %d lines, one per satellite", len(fleet))
		return
	}
	for i, s := range fleet {
		if len(sats) > 0 && !slices.Contains(sats, s.ID) {
			continue
		}
		want := s.ID + " never heard"
		if line, ok := res.world.newest[s.ID]; ok {
			want = line
		}
		if i >= len(lines) {
			t.Errorf("reports/health.txt ends after %d lines, want %d; the next is %q", len(lines), len(fleet), want)
			return
		}
		if lines[i] != want {
			t.Errorf("reports/health.txt line %d = %q, want %q", i+1, lines[i], want)
			return
		}
	}
	if len(sats) == 0 && len(lines) > len(fleet) {
		t.Errorf("reports/health.txt has %d lines, want %d, one per satellite", len(lines), len(fleet))
	}
}

func TestHealth(t *testing.T) {
	c := testThursday("YST", "TIR", "PEL")
	// Ystmark and Port Ellery relay every frame 150 s after it was made: a satellite over
	// Ystmark and Tirrin Rock at once arrives out of number order, and the last Thistle-2
	// frame Port Ellery relays, made at 16:06:56, reads batt_v=nan.
	c.late = map[string]time.Duration{"YST": 150 * time.Second, "PEL": 150 * time.Second}
	res := simRun(t, c)
	if res.err != nil {
		t.Fatalf("rehearse() = %v, want nil", res.err)
	}
	t.Run("Thistle-2's battery sensor", func(t *testing.T) {
		testHealth(t, res, "thistle-2")
	})
	t.Run("out of order at Ystmark and Tirrin Rock", func(t *testing.T) {
		testHealth(t, res)
	})
	t.Run("the investors' week", func(t *testing.T) {
		testHealth(t, simWeekRun(t))
	})
	t.Run("stopped at noon", func(t *testing.T) {
		res := simRun(t, simCase{from: simJuly(26, 0, 0, 0), days: 1, stations: []string{"PEL"}, stopAt: simJuly(26, 12, 0, 0)})
		if _, ok := res.file("reports/health.txt"); ok {
			t.Errorf("reports/health.txt is there after a stop, want none: it comes after a whole run")
		}
	})
}

// testCallsAt returns the calls a station got at an instant, in the order they were made.
func testCallsAt(res *simResult, station string, at time.Time) []*simCall {
	var calls []*simCall
	for _, cl := range res.world.calls {
		if cl.station == station && cl.dialed.Equal(at) {
			calls = append(calls, cl)
		}
	}
	return calls
}

func TestLies(t *testing.T) {
	t.Run("a flipped byte at Varsk Island", func(t *testing.T) {
		c := testThursday("VSK")
		c.corrupt = map[string]int{"cirrine-5/VSK/01487": 3} // its third TLM frame, at 00:44:53
		res := simRun(t, c)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		w := res.world
		if read := w.tlmGood + w.tlmBad; read != w.relayed["VSK"] {
			t.Errorf("rehearsal read %d of the %d TLM frames Varsk Island relayed, want all %d: a bad CRC spoils only its own frame", read, w.relayed["VSK"], w.relayed["VSK"])
		}
		testFlew(t, res, testFlown)
		testHealth(t, res)
	})
	t.Run("a length of 8", func(t *testing.T) {
		c := testThursday("VSK")
		c.badLength = map[string]int{"pipit-2/VSK/22860": 1} // right after its WELCOME, at its AOS
		res := simRun(t, c)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		bad := simJuly(26, 1, 39, 45)
		if calls := testCallsAt(res, "VSK", bad); len(calls) == 0 {
			t.Errorf("no call to VSK at Thu 01:39:45 for pipit-2/VSK/22860, want one at its AOS")
		} else if !calls[0].hungUp.Equal(bad) {
			t.Errorf("call for pipit-2/VSK/22860 hung up at %s, want %s, when the bad length came", testAt(calls[0].hungUp, bad), testAt(bad, calls[0].hungUp))
		} else if len(calls) == 1 {
			t.Errorf("no second call to VSK at Thu 01:39:45 for pipit-2/VSK/22860, want one at once: its line was lost")
		} else if los := simJuly(26, 1, 48, 41); !calls[1].hungUp.Equal(los) {
			t.Errorf("second call for pipit-2/VSK/22860 hung up at %s, want %s, its LOS", testAt(calls[1].hungUp, los), testAt(los, calls[1].hungUp))
		}
		testFlew(t, res, testFlown)
	})
}

// testIdle checks the calls to a station that answers and then says nothing: every one hung
// up 90 s after it was dialed, or at its pass's LOS when that came first, and each pass's
// calls dialed again after the waits testRedials wants. It reports the first problem only.
func testIdle(t *testing.T, res *simResult, station string) {
	t.Helper()
	dialedFor := testDialer(res)
	var passes []Contact
	byPass := map[string][]*simCall{}
	for _, cl := range res.world.calls {
		if cl.station != station || cl.failed {
			continue
		}
		c, ok := dialedFor(cl)
		if !ok {
			continue
		}
		want, why := cl.dialed.Add(90*time.Second), "90 s after it was dialed"
		if c.LOS.Before(want) {
			want, why = c.LOS, "its LOS"
		}
		if !cl.hungUp.Equal(want) {
			t.Errorf("call for %s dialed at %s hung up at %s, want %s, %s", c.ID, cl.dialed.Format("Mon 15:04:05"), testAt(cl.hungUp, want), testAt(want, cl.hungUp), why)
			return
		}
		if byPass[c.ID] == nil {
			passes = append(passes, c)
		}
		byPass[c.ID] = append(byPass[c.ID], cl)
	}
	for _, c := range passes {
		if !testWaits(t, c, byPass[c.ID]) {
			return
		}
	}
}

// testRedials checks the calls made for one contact: after a call the station welcomed
// and then lost, the next comes at once; after the n-th failed call in a row, it comes after
// a wait of at least half of 1 s doubled n-1 times, at most 30 s, and less than all of it.
// It reports the first problem only.
func testRedials(t *testing.T, res *simResult, c Contact) {
	t.Helper()
	dialedFor := testDialer(res)
	var calls []*simCall
	for _, cl := range res.world.calls {
		if cl.station != c.Station {
			continue
		}
		if d, ok := dialedFor(cl); ok && d.ID == c.ID {
			calls = append(calls, cl)
		}
	}
	testWaits(t, c, calls)
}

// testWaits checks the waits between a contact's calls, given in the order they were made,
// as testRedials says. It reports the first problem only, and whether there was none.
func testWaits(t *testing.T, c Contact, calls []*simCall) bool {
	t.Helper()
	failures := 0
	for i := 1; i < len(calls); i++ {
		prev, cl := calls[i-1], calls[i]
		end := prev.hungUp
		if prev.failed {
			end = prev.dialed
		}
		wait := cl.dialed.Sub(end)
		if prev.welcomed {
			failures = 0
			if wait != 0 {
				t.Errorf("call %d for %s came %v after call %d lost its line, want it at once", i+1, c.ID, wait, i)
				return false
			}
			continue
		}
		failures++
		most := 30 * time.Second
		if failures < 6 {
			most = time.Second << (failures - 1)
		}
		if wait < most/2 || wait >= most {
			t.Errorf("call %d for %s came %v after call %d failed, want %v to %v: failure %d in a row", i+1, c.ID, wait, i, most/2, most, failures)
			return false
		}
	}
	return true
}

func TestQuiet(t *testing.T) {
	c := testThursday("YST")
	simStorm(&c)
	res := simRun(t, c)
	if res.err != nil {
		t.Fatalf("rehearse() = %v, want nil", res.err)
	}
	t.Run("the fiber dies at 03:31", func(t *testing.T) {
		p := res.plan().byID["pipit-2/YST/22861"]
		// Its last frame before the fiber died, a TLM frame, came at 03:30:52.
		lost := simJuly(26, 3, 32, 22)
		switch calls := testCallsAt(res, "YST", p.AOS); {
		case len(calls) == 0:
			t.Errorf("no call to YST at Thu 03:28:42 for pipit-2/YST/22861, want one at its AOS")
		case calls[0].farDeadline:
			t.Errorf("call for pipit-2/YST/22861 has a read deadline more than a day away: a deadline is an instant on the time package's clock, never w.Now()'s")
		case !calls[0].hungUp.Equal(lost):
			t.Errorf("call for pipit-2/YST/22861 hung up at %s, want Thu 03:32:22, 90 s after its last frame at 03:30:52", testAt(calls[0].hungUp, lost))
		default:
			testRedials(t, res, p)
		}
		pl := testByPass(res)[p.ID]
		if pl == nil {
			pl = &testPassLines{}
		}
		outcome := testCommanded(res, func(Contact) string { return "failed line_lost" })
		testOne(t, "LOS", p.ID, pl.los, p.LOS, outcome(p))
	})
	t.Run("the backhaul from 03:41", func(t *testing.T) {
		// Until 19:40 every message takes 0.35 s longer each way, and Ystmark's own arrive
		// in pieces of at most 64 bytes, 0.25 s apart.
		up, outcome, dialedFor := simJuly(26, 3, 41, 0), testCommanded(res, testWeather(c)), testDialer(res)
		calls := map[string]int{}
		for _, cl := range res.world.calls {
			if ct, ok := dialedFor(cl); ok {
				calls[ct.ID]++
			}
		}
		byPass := testByPass(res)
		for _, ct := range res.plan().contacts {
			if ct.AOS.Before(up) {
				continue
			}
			pl := byPass[ct.ID]
			if pl == nil {
				pl = &testPassLines{}
			}
			if !testOne(t, "AOS", ct.ID, pl.aos, ct.AOS, fmt.Sprintf("dish %d", ct.Dish)) || !testOne(t, "LOS", ct.ID, pl.los, ct.LOS, outcome(ct)) {
				return
			}
			if calls[ct.ID] != 1 {
				t.Errorf("rehearsal made %d calls for %s, want 1: its line is slow, not dead", calls[ct.ID], ct.ID)
				return
			}
		}
	})
	t.Run("stopped while dialing again", func(t *testing.T) {
		c := c
		stop := simJuly(26, 3, 33, 0)
		c.stopAt = stop
		res := simRun(t, c)
		if !errors.Is(res.err, context.Canceled) {
			t.Errorf("rehearse() = %v, want context.Canceled", res.err)
		}
		for _, cl := range res.world.calls {
			if cl.dialed.After(stop) || cl.late {
				t.Errorf("call to YST dialed at %s, want none after the stop at Thu 03:33:00", testAt(cl.dialed, stop))
				break
			}
		}
		p := res.plan().byID["pipit-2/YST/22861"]
		if pl := testByPass(res)[p.ID]; pl == nil {
			t.Errorf("console has no LOS line for pipit-2/YST/22861, want %q", "Thu 03:33:00 LOS pipit-2/YST/22861 stopped")
		} else {
			testOne(t, "LOS", p.ID, pl.los, stop, "stopped")
		}
		want, last := "rehearsal stopped at Thu 03:33:00", ""
		if lines := res.lines(); len(lines) > 0 {
			last = lines[len(lines)-1]
		}
		if last != want {
			t.Errorf("last console line = %q, want %q", last, want)
		} else if !res.ended.Equal(stop) {
			t.Errorf("rehearse() returned at %s, want Thu 03:33:00, at the stop", testAt(res.ended, stop))
		}
	})
}

// testCirrine reports whether a satellite is one of Cirrine Atmospherics' thirteen.
func testCirrine(satID string) bool {
	for _, s := range simNet().Satellites {
		if s.ID == satID {
			return s.Operator == "Cirrine Atmospherics"
		}
	}
	return false
}

// testGiven returns the commands a satellite was given from one pass, in the order it got
// them.
func testGiven(res *simResult, ct Contact) []simCommand {
	var out []simCommand
	if s := res.world.sats[ct.SatID]; s != nil {
		for _, cmd := range s.commands {
			if cmd.PassID == ct.ID {
				out = append(out, cmd)
			}
		}
	}
	return out
}

// testShort shortens a command line's long hex fields to their first 8 digits and last 4,
// around "…".
func testShort(line string) string {
	f := strings.Fields(line)
	for i, w := range f {
		prefix, hex := "", w
		if strings.HasPrefix(w, "sig=") {
			prefix, hex = "sig=", w[4:]
		}
		if len(hex) > 16 && strings.Trim(hex, "0123456789abcdef") == "" {
			f[i] = prefix + hex[:8] + "…" + hex[len(hex)-4:]
		}
	}
	return strings.Join(f, " ")
}

// testRan checks the commands each satellite other than Cirrine's ran in a run (TestSealed
// checks theirs): every one it was given ran, numbered one above the last it ran before, the
// first above the last it had at the window's start. It reports the first problem only.
func testRan(t *testing.T, res *simResult) {
	t.Helper()
	for _, s := range simNet().Satellites {
		sat := res.world.sats[s.ID]
		if sat == nil || testCirrine(s.ID) {
			continue
		}
		last, ok := res.c.lastSeq[s.ID]
		if !ok {
			last = simFirstLast(s.ID)
		}
		for _, cmd := range sat.commands {
			switch {
			case cmd.Refused != "":
				t.Errorf("%s refused %s %d from %s at %s: %s", s.ID, cmd.Op, cmd.Seq, cmd.PassID, cmd.At.Format("Mon 15:04:05"), cmd.Refused)
			case cmd.Seq != last+1:
				t.Errorf("%s ran %s %d from %s, want it numbered %d, one above the last it ran", s.ID, cmd.Op, cmd.Seq, cmd.PassID, last+1)
			default:
				last = cmd.Seq
				continue
			}
			return
		}
	}
}

// testTuesday is TestCommands' run: Tuesday at Port Ellery, before k-2029-07 activates, so
// every satellite holds k-2028-10 alone, with Thistle-2 starting from last 0. It is flown
// once per test binary.
func testTuesday(t *testing.T) *simResult {
	t.Helper()
	c := simCase{from: simJuly(24, 0, 0, 0), days: 1, stations: []string{"PEL"}}
	c.lastSeq = map[string]uint64{"thistle-2": 0}
	return simOnce(t, "Tuesday at Port Ellery", c)
}

func TestCommands(t *testing.T) {
	res := testTuesday(t)
	if res.err != nil {
		t.Fatalf("rehearse() = %v, want nil", res.err)
	}
	t.Run("a PING on every pass", func(t *testing.T) {
		for _, ct := range res.plan().contacts {
			if testCirrine(ct.SatID) {
				continue
			}
			cmds, at := testGiven(res, ct), ct.AOS.Add(12*time.Second)
			switch {
			case len(cmds) == 0:
				t.Errorf("%s: its satellite was given no command, want a PING sent at %s, 10 s after AOS", ct.ID, ct.AOS.Add(10*time.Second).Format("Mon 15:04:05"))
			case cmds[0].Op != "PING":
				t.Errorf("%s: its first command is %q, want a PING", ct.ID, testShort(cmds[0].Line))
			case !cmds[0].At.Equal(at):
				t.Errorf("%s's PING ran at %s, want %s: sent 10 s after AOS and run 2 s later", ct.ID, testAt(cmds[0].At, at), testAt(at, cmds[0].At))
			default:
				continue
			}
			return
		}
		testRan(t, res)
		testFlew(t, res, testFlown)
	})
	t.Run("Pipit-2's first line", func(t *testing.T) {
		want := "CMD pipit-2 6315 PING sig=0cb110af1256712defcb30409a53e47e402b944c6a13ab67a236315160034810"
		if s := res.world.sats["pipit-2"]; s == nil || len(s.commands) == 0 {
			t.Errorf("pipit-2 was given no command, want its first line %q", testShort(want))
		} else if got := s.commands[0].Line; got != want {
			t.Errorf("pipit-2's first line = %q, want %q", testShort(got), testShort(want))
		}
	})
	t.Run("Thistle-2 from 0", func(t *testing.T) {
		switch s := res.world.sats["thistle-2"]; {
		case s == nil || len(s.commands) == 0:
			t.Errorf("thistle-2 was given no command, want a PING numbered 1: its WELCOME said last 0")
		case s.commands[0].Seq != 1:
			t.Errorf("thistle-2's first command = %q, want it numbered 1: its WELCOME said last 0", testShort(s.commands[0].Line))
		case s.commands[0].Refused != "":
			t.Errorf("thistle-2 refused its first command, %q: %s", testShort(s.commands[0].Line), s.commands[0].Refused)
		}
	})
}

// testOneAtATime checks the commands every satellite was given in a run, in the order it
// got them: none refused as replayed, except a copy sent again of one it had already run,
// and each answered at least 2 s after the one before. A satellite answers a command 2 s
// after it arrives, so two closer than that were in flight at once. It reports the first
// problem only.
func testOneAtATime(t *testing.T, res *simResult) {
	t.Helper()
	for _, s := range simNet().Satellites {
		sat := res.world.sats[s.ID]
		if sat == nil {
			continue
		}
		ran := map[uint64]string{} // the pass each number ran on
		for i, cmd := range sat.commands {
			if strings.HasPrefix(cmd.Refused, "replayed sequence number") && ran[cmd.Seq] != cmd.PassID {
				t.Errorf("%s refused %s %d from %s: %s; want it above every number sent before", s.ID, cmd.Op, cmd.Seq, cmd.PassID, cmd.Refused)
				return
			}
			if i > 0 {
				prev := sat.commands[i-1]
				if gap := cmd.At.Sub(prev.At); gap < 2*time.Second {
					t.Errorf("%s answered %s %d (%s) %v after %s %d (%s), want one in flight at a time", s.ID, cmd.Op, cmd.Seq, cmd.PassID, gap, prev.Op, prev.Seq, prev.PassID)
					return
				}
			}
			if cmd.Refused == "" {
				ran[cmd.Seq] = cmd.PassID
			}
		}
	}
}

// testSentAll checks that every contact of a run whose first call the station welcomed
// gave its satellite a command. It reports the first problem only.
func testSentAll(t *testing.T, res *simResult) {
	t.Helper()
	sent := map[string]bool{} // the pass IDs a command came from
	for _, s := range res.world.sats {
		for _, cmd := range s.commands {
			sent[cmd.PassID] = true
		}
	}
	plan := res.plan()
	for _, cl := range res.world.calls {
		if c, ok := plan.byID[cl.passID]; ok && cl.welcomed && cl.dialed.Equal(c.AOS) && !sent[c.ID] {
			t.Errorf("%s: its satellite was given no command, want a PING once its lock is free", c.ID)
			return
		}
	}
}

func TestNumbers(t *testing.T) {
	t.Run("Ystmark and Tirrin Rock", func(t *testing.T) {
		// Over the pole the two stations often hold one satellite at once: Furrow-O
		// rises over Ystmark at 18:26:27 and over Tirrin Rock 2 s later.
		res := simRun(t, testThursday("YST", "TIR"))
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testOneAtATime(t, res)
	})
	t.Run("a line dead until LOS", func(t *testing.T) {
		// Port Ellery's line dies a second after Thistle-2's 16:45:15 PING goes up, and
		// stays dead past the pass's LOS: the PING runs, and its answer never comes.
		c := simCase{from: simJuly(24, 0, 0, 0), days: 1, stations: []string{"PEL"}}
		c.passes = testArchive(c, "thistle-2/PEL/02106", "thistle-2/PEL/02107")
		c.quiet = []simSpan{{Station: "PEL", From: simJuly(24, 16, 45, 26), To: simJuly(24, 17, 0, 0)}}
		res := simRun(t, c)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		plan, byPass := res.plan(), testByPass(res)
		for _, id := range []string{"thistle-2/PEL/02106", "thistle-2/PEL/02107"} {
			ct, outcome := plan.byID[id], "flown"
			if id == "thistle-2/PEL/02106" {
				outcome = "failed line_lost"
			} else if len(testGiven(res, ct)) == 0 {
				t.Errorf("%s: thistle-2 was given no command, want a PING: the pass before it ended at its LOS", id)
				return
			}
			pl := byPass[id]
			if pl == nil {
				pl = &testPassLines{}
			}
			if !testOne(t, "LOS", id, pl.los, ct.LOS, outcome) {
				return
			}
		}
	})
	t.Run("the investors' week", func(t *testing.T) {
		res := simWeekRun(t)
		testOneAtATime(t, res)
		testSentAll(t, res)
	})
}

// testSealPING returns the SEAL line of a PING to cirrine-1 numbered 9713 under k-2028-10,
// sealed with the nonce given.
func testSealPING(nonce []byte) string {
	key := simSatKey([]byte("uplink next master key, 32 bytes"), "seal", "cirrine-1")
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	sealed := gcm.Seal(nil, nonce, []byte("PING"), []byte("SEAL|2|cirrine-1|9713|k-2028-10"))
	return fmt.Sprintf("SEAL 2 cirrine-1 9713 k-2028-10 %x %x", nonce, sealed)
}

func TestSealed(t *testing.T) {
	res := simWeekRun(t)
	t.Run("Cirrine's thirteen", func(t *testing.T) {
		for _, s := range simNet().Satellites {
			sat := res.world.sats[s.ID]
			if sat == nil {
				continue
			}
			cirrine := testCirrine(s.ID)
			for _, cmd := range sat.commands {
				switch {
				case cirrine && !cmd.Sealed:
					t.Errorf("%s was given %q from %s, want a SEAL line: Cirrine's commands are all sealed", s.ID, testShort(cmd.Line), cmd.PassID)
				case cirrine && cmd.Refused != "":
					t.Errorf("%s refused %q from %s: %s", s.ID, testShort(cmd.Line), cmd.PassID, cmd.Refused)
				case !cirrine && cmd.Op == "PING" && cmd.Sealed:
					t.Errorf("%s was given %q from %s, want a CMD line: only Cirrine's PINGs are sealed", s.ID, testShort(cmd.Line), cmd.PassID)
				default:
					continue
				}
				return
			}
		}
	})
	t.Run("Cirrine-1's first line", func(t *testing.T) {
		fixed, _ := hex.DecodeString("a1b2c3d4e5f60718293a4b5c")
		want := "SEAL 2 cirrine-1 9713 k-2028-10 a1b2c3d4e5f60718293a4b5c fc5a75921a95c7faaf06adb91ebc66ea2d2d379c"
		if got := testSealPING(fixed); got != want {
			t.Fatalf("the test seals %q, want %q", got, want)
		}
		s := res.world.sats["cirrine-1"]
		if s == nil || len(s.commands) == 0 {
			t.Errorf("cirrine-1 was given no command, want its first line like %q", testShort(want))
			return
		}
		got := s.commands[0].Line
		f := strings.Fields(got)
		var nonce []byte
		if len(f) == 7 {
			nonce, _ = hex.DecodeString(f[5])
		}
		if len(nonce) != 12 || !strings.HasPrefix(got, "SEAL 2 cirrine-1 9713 k-2028-10 ") {
			t.Errorf("cirrine-1's first line = %q, want %q with a nonce of its own", testShort(got), testShort(want))
		} else if w := testSealPING(nonce); got != w {
			t.Errorf("cirrine-1's first line = %q, want %q: PING under its seal key, with SEAL|2|cirrine-1|9713|k-2028-10 as associated data", testShort(got), testShort(w))
		}
	})
	t.Run("a fresh nonce every time", func(t *testing.T) {
		seen := map[string]string{} // each nonce, and the line it sealed
		for _, run := range []*simResult{res, testTuesday(t)} {
			for _, s := range simNet().Satellites {
				sat := run.world.sats[s.ID]
				if sat == nil {
					continue
				}
				for _, cmd := range sat.commands {
					f := strings.Fields(cmd.Line)
					if !cmd.Sealed || len(f) != 7 {
						continue
					}
					if line, ok := seen[f[5]]; ok && line != cmd.Line {
						t.Errorf("nonce %s sealed %q and %q, want a fresh one for every seal", testShort(f[5]), testShort(line), testShort(cmd.Line))
						return
					}
					seen[f[5]] = cmd.Line
				}
			}
		}
		// With crypto/rand's stream fixed by a seed, every nonce is a piece of it.
		cryptotest.SetGlobalRandom(t, 7)
		stream := make([]byte, 1<<16)
		rand.Read(stream)
		cryptotest.SetGlobalRandom(t, 7)
		ids := []string{"cirrine-6/PEL/01458", "cirrine-1/PEL/20255"}
		c := simCase{from: simJuly(24, 0, 0, 0), days: 1, stations: []string{"PEL"}}
		c.passes = testArchive(c, ids...)
		fixed := simRun(t, c)
		for _, id := range ids {
			sealed := 0
			for _, cmd := range testGiven(fixed, fixed.plan().byID[id]) {
				f := strings.Fields(cmd.Line)
				if !cmd.Sealed || len(f) != 7 {
					continue
				}
				sealed++
				if nonce, err := hex.DecodeString(f[5]); err != nil || !bytes.Contains(stream, nonce) {
					t.Errorf("%s: nonce %s doesn't come from crypto/rand, want every nonce from it: math/rand's can be predicted", id, testShort(f[5]))
					return
				}
			}
			if sealed == 0 {
				t.Errorf("%s: its satellite was given no SEAL line, want its PING sealed with a nonce from crypto/rand", id)
				return
			}
		}
	})
}

// testUplinks returns the UPLINK frames a station received on a call, in order.
func testUplinks(cl *simCall) []simFrame {
	var out []simFrame
	br := bufio.NewReader(bytes.NewReader(cl.received))
	for {
		f, err := simReadFrame(br)
		if err != nil {
			return out
		}
		if f.Type == TypeUplink {
			out = append(out, f)
		}
	}
}

// testPings returns the PINGs a satellite was given from one pass, in order.
func testPings(res *simResult, ct Contact) []simCommand {
	var out []simCommand
	for _, cmd := range testGiven(res, ct) {
		if cmd.Op == "PING" {
			out = append(out, cmd)
		}
	}
	return out
}

// testDropped flies Thursday at Ystmark with one pass, its trouble set by trouble, and
// checks that the pass ends flown.
func testDropped(t *testing.T, id string, trouble func(c *simCase)) (*simResult, Contact) {
	t.Helper()
	c := testThursday("YST")
	c.passes = testArchive(c, id)
	trouble(&c)
	res := simRun(t, c)
	if res.err != nil {
		t.Fatalf("rehearse() = %v, want nil", res.err)
	}
	ct := res.plan().byID[id]
	pl := testByPass(res)[id]
	if pl == nil {
		pl = &testPassLines{}
	}
	testOne(t, "LOS", id, pl.los, ct.LOS, "flown")
	return res, ct
}

func TestNeverTwice(t *testing.T) {
	const id = "cirrine-7/YST/01488" // 01:51:00 to 01:57:56
	t.Run("dropped after it ran", func(t *testing.T) {
		res, ct := testDropped(t, id, func(c *simCase) { c.dropAfter = map[string]int{id: 1} })
		if cmds := testGiven(res, ct); len(cmds) != 1 {
			t.Errorf("cirrine-7 was given %d commands from %s, want 1: the next WELCOME's last says its PING ran", len(cmds), id)
		}
	})
	t.Run("dropped before it ran", func(t *testing.T) {
		res, ct := testDropped(t, id, func(c *simCase) { c.dropBefore = map[string]int{id: 1} })
		var first []simFrame
		if calls := testCallsAt(res, "YST", ct.AOS); len(calls) > 0 {
			first = testUplinks(calls[0])
		}
		switch pings := testPings(res, ct); {
		case len(first) == 0:
			t.Errorf("the call for %s at its AOS carried no command, want its PING", id)
		case len(pings) == 0:
			t.Errorf("cirrine-7 was given no PING from %s, want the one the dropped call carried, sent again", id)
		case pings[0].Line != string(first[0].Body):
			t.Errorf("cirrine-7 was given %q, want the dropped call's line again: %q", testShort(pings[0].Line), testShort(string(first[0].Body)))
		case len(pings) > 1:
			t.Errorf("cirrine-7 was given %d PINGs from %s, want 1", len(pings), id)
		}
	})
	t.Run("lost on its way", func(t *testing.T) {
		// Over the backhaul, cirrine-1's PING of 04:25:12 reaches Ystmark at 04:25:12.35,
		// and the bad length after the next frame loses the line before it runs at 04:25:14.35.
		const id = "cirrine-1/YST/20282"
		res, ct := testDropped(t, id, func(c *simCase) {
			simStorm(c)
			c.badLength = map[string]int{id: 2}
		})
		ran := 0
		for _, cmd := range testPings(res, ct) {
			if cmd.Refused == "" {
				ran++
			}
		}
		if ran != 1 {
			t.Errorf("cirrine-1 ran %d PINGs from %s, want 1", ran, id)
		}
	})
}

// testKeyIDs returns the operation and key ID of each of a satellite's first commands in a
// run, up to n of them: "PING k-2028-10".
func testKeyIDs(res *simResult, satID string, n int) []string {
	var out []string
	if s := res.world.sats[satID]; s != nil {
		for _, cmd := range s.commands[:min(n, len(s.commands))] {
			out = append(out, cmd.Op+" "+cmd.KeyID)
		}
	}
	return out
}

func TestKeys(t *testing.T) {
	t.Run("the investors' week", func(t *testing.T) {
		// Every satellite holds k-2028-10 alone on Monday at 00:00, and k-2029-07
		// activates on Wednesday at 00:00.
		res := simWeekRun(t)
		switchAt := simJuly(25, 0, 0, 0)
		loaded, last := 0, time.Time{}
		for _, s := range simNet().Satellites {
			sat := res.world.sats[s.ID]
			if sat == nil {
				continue
			}
			for i, cmd := range sat.commands {
				want := "k-2028-10"
				if !cmd.At.Before(switchAt) && cmd.Op != "LOAD_KEY" {
					want = "k-2029-07"
				}
				load := cmd.Op == "LOAD_KEY" && len(cmd.Args) > 0 && cmd.Args[0] == "k-2029-07"
				switch {
				case cmd.Refused != "":
					t.Errorf("%s refused %q from %s: %s", s.ID, testShort(cmd.Line), cmd.PassID, cmd.Refused)
				case i == 1 && !load:
					t.Errorf("%s's second command, from %s, is %q, want LOAD_KEY for k-2029-07", s.ID, cmd.PassID, testShort(cmd.Line))
				case i != 1 && cmd.Op == "LOAD_KEY":
					t.Errorf("%s was given LOAD_KEY as its command %d, from %s, want it second only, after its first PING", s.ID, i+1, cmd.PassID)
				case cmd.KeyID != want:
					t.Errorf("%s's %s %d at %s came under %s, want %s", s.ID, cmd.Op, cmd.Seq, cmd.At.Format("Mon 15:04:05"), cmd.KeyID, want)
				default:
					if load {
						loaded++
						last = cmd.At
					}
					continue
				}
				return
			}
		}
		// Thistle-2 is the last to rise, at Lantern Atoll at 09:27:17.
		if by := simJuly(23, 9, 27, 31); !raceOn && (loaded != 80 || last.After(by)) {
			t.Errorf("%d satellites loaded k-2029-07, the last at %s, want all 80 by Mon 09:27:31", loaded, testAt(last, by))
		}
	})
	t.Run("from Wednesday at Port Ellery", func(t *testing.T) {
		// A window that starts as k-2029-07 activates: no satellite holds it yet.
		res := simRun(t, simCase{from: simJuly(25, 0, 0, 0), days: 1, stations: []string{"PEL"}})
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		want := []string{"PING k-2028-10", "LOAD_KEY k-2028-10", "PING k-2029-07"}
		for _, s := range simNet().Satellites {
			if sat := res.world.sats[s.ID]; sat == nil || len(sat.commands) < 3 {
				continue
			}
			for i, cmd := range res.world.sats[s.ID].commands[:3] {
				if cmd.Refused != "" {
					t.Errorf("%s refused %q from %s: %s", s.ID, testShort(cmd.Line), cmd.PassID, cmd.Refused)
					return
				}
				if got := cmd.Op + " " + cmd.KeyID; got != want[i] {
					t.Errorf("%s's first commands = %q, want %q", s.ID, testKeyIDs(res, s.ID, 3), want)
					return
				}
			}
		}
	})
	t.Run("a LOAD_KEY unanswered at LOS", func(t *testing.T) {
		// Port Ellery's line dies a second after Thistle-2's 16:45:15 LOAD_KEY goes up,
		// and stays dead past the pass's LOS: the LOAD_KEY runs, and its answer never comes.
		c := simCase{from: simJuly(24, 0, 0, 0), days: 1, stations: []string{"PEL"}}
		c.passes = testArchive(c, "thistle-2/PEL/02106", "thistle-2/PEL/02107")
		c.quiet = []simSpan{{Station: "PEL", From: simJuly(24, 16, 45, 28), To: simJuly(24, 17, 0, 0)}}
		res := simRun(t, c)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		plan, byPass := res.plan(), testByPass(res)
		for _, id := range []string{"thistle-2/PEL/02106", "thistle-2/PEL/02107"} {
			ct, outcome := plan.byID[id], "flown"
			if id == "thistle-2/PEL/02106" {
				outcome = "failed line_lost"
			}
			loads := 0
			for _, cmd := range testGiven(res, ct) {
				if cmd.Op == "LOAD_KEY" {
					loads++
				}
			}
			switch {
			case id == "thistle-2/PEL/02106" && loads != 1:
				t.Errorf("thistle-2 was given %d LOAD_KEYs from %s, want 1, after its PING", loads, id)
			case id == "thistle-2/PEL/02107" && loads != 0:
				t.Errorf("thistle-2 was given LOAD_KEY again from %s, want none: the WELCOME's last says the first ran", id)
			default:
				pl := byPass[id]
				if pl == nil {
					pl = &testPassLines{}
				}
				if testOne(t, "LOS", id, pl.los, ct.LOS, outcome) {
					continue
				}
			}
			return
		}
	})
}

// testBooked returns the week's bookings of the passes given, from a case's archive.
func testBooked(c simCase, ids ...string) []Booking {
	var out []Booking
	for _, b := range simWeekBookings(simWindow(c.from, c.days), c.stations) {
		for _, p := range simArchive(simWindow(c.from, c.days), c.stations) {
			if slices.Contains(ids, p.ID) && p.SatID == b.SatID && p.Station == b.Station && p.AOS.Equal(b.AOS) {
				out = append(out, b)
			}
		}
	}
	return out
}

// testLandings returns how many LANDED lines the console has for each file.
func testLandings(res *simResult) map[string]int {
	out := map[string]int{}
	for _, l := range res.console {
		if f := strings.Fields(l.Text); len(f) == 4 && f[2] == "LANDED" {
			out[f[3]]++
		}
	}
	return out
}

// testWhy says how a file's chunks reached the rehearsal on its first pass, when one came
// after a higher one or twice.
func testWhy(res *simResult, name string) string {
	first, last := "", -1
	seen := map[int]bool{}
	for _, ch := range res.world.chunks {
		if ch.File != name || (first != "" && ch.PassID != first) {
			continue
		}
		first = ch.PassID
		switch {
		case seen[ch.N]:
			return fmt.Sprintf(": chunk %d came twice", ch.N)
		case ch.N < last:
			return fmt.Sprintf(": chunk %d came after chunk %d", ch.N, last)
		}
		seen[ch.N], last = true, ch.N
	}
	return ""
}

// testFile returns what is wrong with a file a satellite sent down, or "". With why ""
// it must be in files/, byte for byte as sent, with one LANDED line; otherwise, being
// incomplete for that reason, it must not be there and have no LANDED line.
func testFile(res *simResult, landings map[string]int, name, why string) string {
	sent := res.world.files[name]
	got, ok := res.file("files/" + name)
	switch {
	case sent == nil:
		return fmt.Sprintf("no satellite sent %s, want it sent: the test's case is wrong", name)
	case why == "?" && (ok || landings[name] > 0):
		return fmt.Sprintf("%s landed, want it left out: none of its chunks came", name)
	case why != "" && ok:
		return fmt.Sprintf("files/%s is there, want it left out: %s", name, why)
	case why != "" && landings[name] > 0:
		return fmt.Sprintf("console says LANDED %s, want no such line: %s", name, why)
	case why != "":
		return ""
	case !ok:
		return fmt.Sprintf("files/%s is not there, want the %d bytes %s sent", name, len(sent.data), sent.booking.SatID)
	case len(got) != len(sent.data):
		return fmt.Sprintf("files/%s holds %d bytes, want the %d %s sent%s", name, len(got), len(sent.data), sent.booking.SatID, testWhy(res, name))
	case got != string(sent.data):
		return fmt.Sprintf("files/%s is not the file %s sent%s", name, sent.booking.SatID, testWhy(res, name))
	case landings[name] != 1:
		return fmt.Sprintf("console has %d LANDED lines for %s, want one", landings[name], name)
	}
	return ""
}

// testIncomplete returns why each file a run's satellites sent down can't be whole, or ""
// when it can: "missing 7 8" when chunks of it never reached the rehearsal, "crc mismatch"
// when its END carried a wrong CRC and no later pass sent it again, "?" when none of its
// chunks did, so the rehearsal may never have heard of it.
func testIncomplete(res *simResult) map[string]string {
	arrived := map[string]map[int]bool{}   // each file's chunk numbers that came
	passes := map[string]map[string]bool{} // the passes they came on
	for _, ch := range res.world.chunks {
		if arrived[ch.File] == nil {
			arrived[ch.File], passes[ch.File] = map[int]bool{}, map[string]bool{}
		}
		arrived[ch.File][ch.N], passes[ch.File][ch.PassID] = true, true
	}
	out := map[string]string{}
	for name, f := range res.world.files {
		var missing []string
		for n := range len(f.chunks) {
			if !arrived[name][n] {
				missing = append(missing, strconv.Itoa(n))
			}
		}
		switch {
		case len(arrived[name]) == 0:
			out[name] = "?"
		case len(missing) > 0:
			out[name] = "missing " + strings.Join(missing, " ")
		case slices.Contains(res.c.badCRC, name) && len(passes[name]) < 2:
			out[name] = "crc mismatch"
		default:
			out[name] = ""
		}
	}
	return out
}

// testFiles checks the files a run's satellites sent down, the ones named or all of them:
// each whole in files/, byte for byte, with one LANDED line, unless testIncomplete says why
// it can't be; nothing else in files/; and reports/incomplete.txt listing the incomplete
// ones by name, with their reasons. It reports the first problem of each kind.
func testFiles(t *testing.T, res *simResult, only ...string) {
	t.Helper()
	if res.err != nil {
		t.Fatalf("rehearse() = %v, want nil", res.err)
	}
	incomplete := testIncomplete(res)
	var names, want []string
	for name := range incomplete {
		names = append(names, name)
	}
	sort.Strings(names)
	landings := testLandings(res)
	first, more := "", 0
	for _, name := range names {
		why := incomplete[name]
		if why != "" && why != "?" {
			want = append(want, name+" "+why)
		}
		if len(only) > 0 && !slices.Contains(only, name) {
			continue
		}
		if msg := testFile(res, landings, name, why); msg != "" && first == "" {
			first = msg
		} else if msg != "" {
			more++
		}
	}
	switch {
	case more > 0:
		t.Errorf("%s; and %d more files are wrong", first, more)
	case first != "":
		t.Error(first)
	}
	for _, name := range res.names("files") {
		if _, ok := res.world.files[name]; !ok {
			t.Errorf("files/%s is there, want only the files satellites sent", name)
			break
		}
	}
	lines, ok := res.fileLines("reports/incomplete.txt")
	lines = slices.DeleteFunc(lines, func(l string) bool { name, _, _ := strings.Cut(l, " "); return incomplete[name] == "?" })
	switch {
	case !ok && len(want) == 0:
		t.Errorf("reports/incomplete.txt is not there, want it there and empty")
	case !ok:
		t.Errorf("reports/incomplete.txt is not there, want %q", want)
	case !slices.Equal(lines, want):
		t.Errorf("reports/incomplete.txt = %q, want %q", lines, want)
	}
}

// testLandedAt checks a file's LANDED line: written at the instant at, its END's.
func testLandedAt(t *testing.T, res *simResult, name string, at time.Time) {
	t.Helper()
	want := at.Format("Mon 15:04:05") + " LANDED " + name
	got := res.find(" LANDED " + name + "\n")
	switch {
	case len(got) == 0:
		t.Errorf("console has no LANDED line for %s, want %q", name, want)
	case !got[0].At.Equal(at):
		t.Errorf("LANDED line for %s at %s, want %s, when its END came", name, testAt(got[0].At, at), testAt(at, got[0].At))
	case got[0].Text != want+"\n":
		t.Errorf("LANDED line = %q, want %q", strings.TrimSuffix(got[0].Text, "\n"), want)
	}
}

func TestFiles(t *testing.T) {
	t.Run("the investors' week", func(t *testing.T) {
		// On Tuesday Furrow-K's 07:40 image over Vorrel Island never sends chunks 21 to 23
		// of its own: "furrow-k-20290724-0740.tif missing 21 22 23" until they are asked for.
		res := simWeekRun(t)
		testFiles(t, res)
		if t.Failed() {
			return
		}
		testLandedAt(t, res, "thistle-2-20290726-1603.dat", simJuly(26, 16, 4, 25))
		if !raceOn {
			testLandedAt(t, res, "cirrine-5-20290724-0931.dat", simJuly(24, 9, 32, 18))
		}
	})
	// Thursday at Port Ellery with three booked passes: Cirrine-8's line is lost after the
	// sixth chunk of its file, Thistle-2's file ends with a wrong CRC, and Brine-3's never
	// sends chunks 7 and 8.
	c := testThursday("PEL")
	ids := []string{"cirrine-8/PEL/01496", "thistle-2/PEL/02136", "brine-3/PEL/18374"}
	c.passes, c.bookings = testArchive(c, ids...), testBooked(c, ids...)
	c.badLength = map[string]int{"cirrine-8/PEL/01496": 10}
	c.badCRC = []string{"thistle-2-20290726-1603.dat"}
	c.missing = map[string][]int{"brine-3-20290726-1706.dat": {7, 8}}
	res := simRun(t, c)
	t.Run("a line lost mid-file", func(t *testing.T) {
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		if msg := testFile(res, testLandings(res), "cirrine-8-20290726-1447.dat", ""); msg != "" {
			t.Errorf("%s; its line was lost after chunk 5, and the next call sent FILE again", msg)
		}
	})
	t.Run("a wrong CRC, missing chunks", func(t *testing.T) {
		// No later pass of either satellite: "brine-3-20290726-1706.dat missing 7 8" and
		// "thistle-2-20290726-1603.dat crc mismatch".
		testFiles(t, res, "brine-3-20290726-1706.dat", "thistle-2-20290726-1603.dat")
	})
}

// testEventFields are an event's JSON fields, in the order json.Marshal writes them.
var testEventFields = []string{"id", "type", "booking_id", "sat_id", "station", "file", "bytes", "landed_at"}

// testEventDiff names the first field where an event's body differs from want, as
// "station = \"TIR\", want \"YST\"", or returns "".
func testEventDiff(body string, want Event) string {
	wantBody, _ := json.Marshal(want)
	var got, exp map[string]any
	json.Unmarshal([]byte(body), &got)
	json.Unmarshal(wantBody, &exp)
	for _, k := range testEventFields {
		if fmt.Sprint(got[k]) != fmt.Sprint(exp[k]) {
			return fmt.Sprintf("%s = %s, want %s", k, testJSON(got[k]), testJSON(exp[k]))
		}
	}
	if body != string(wantBody) {
		return "its body isn't the Event as json.Marshal writes it"
	}
	return ""
}

// testJSON writes a decoded JSON value as JSON, or "nothing" when it is missing.
func testJSON(v any) string {
	if v == nil {
		return "nothing"
	}
	data, _ := json.Marshal(v)
	return string(data)
}

// testAnswered reports whether a receiver took a request: it answered 2xx.
func testAnswered(a simAttempt) bool { return a.Status/100 == 2 }

// testRequests checks every request the receivers of a run got: sent to a receiver, not
// refused by it (only planted trouble may fail one), its Apogee-Timestamp the Unix second it
// was sent, its Apogee-Signature 64 lowercase hex digits, as application/json. It reports the
// first problem only.
func testRequests(t *testing.T, res *simResult) {
	t.Helper()
	if len(res.world.hooks.attempts) == 0 && len(testLandings(res)) > 0 {
		t.Errorf("the receivers got no request, want a signed POST for each file that landed")
		return
	}
	for _, a := range res.world.hooks.attempts {
		ts := strconv.FormatInt(a.At.Unix(), 10)
		switch {
		case a.Customer == "":
			t.Errorf("%s for %s went to no customer's receiver: %s; want its customer's webhook_url", a.Event.ID, a.Event.File, a.Problem)
		case a.Problem != "" && a.Problem != "receiver in trouble":
			t.Errorf("%s's receiver refused %s for %s: %s", a.Customer, a.Event.ID, a.Event.File, a.Problem)
		case a.Timestamp != ts:
			t.Errorf("%s for %s sent at %s with Apogee-Timestamp %q, want %q, that instant's Unix second", a.Event.ID, a.Event.File, a.At.Format("Mon 15:04:05"), a.Timestamp, ts)
		case len(a.Signature) != 64 || strings.Trim(a.Signature, "0123456789abcdef") != "":
			t.Errorf("%s for %s came with Apogee-Signature %q, want 64 lowercase hex digits", a.Event.ID, a.Event.File, a.Signature)
		case a.Type != "application/json":
			t.Errorf("%s for %s came with Content-Type %q, want \"application/json\"", a.Event.ID, a.Event.File, a.Type)
		default:
			continue
		}
		return
	}
}

// testEvents checks the events of a run: one for each booked file that landed, first sent
// as it landed, to the customer named for its satellite's operator, saying what landed,
// with IDs from evt-000001 in landing order; a DELIVERED line for each event a receiver
// took, as it took it; and, unless the run was stopped, an UNDELIVERED line for each
// event none took, at its last attempt. It reports the first problem only.
func testEvents(t *testing.T, res *simResult) {
	t.Helper()
	landedAt := map[string]time.Time{}
	var order []string // the files, as they landed
	for _, l := range res.console {
		if f := strings.Fields(l.Text); len(f) == 4 && f[2] == "LANDED" {
			landedAt[f[3]] = l.At
			order = append(order, f[3])
		}
	}
	byFile := map[string][]simAttempt{}
	for _, a := range res.world.hooks.attempts {
		byFile[a.Event.File] = append(byFile[a.Event.File], a)
	}
	delivered, undelivered := map[string][]simLine{}, map[string][]simLine{}
	for _, l := range res.console {
		if f := strings.Fields(l.Text); len(f) == 4 && f[2] == "DELIVERED" {
			delivered[f[3]] = append(delivered[f[3]], l)
		} else if len(f) == 4 && f[2] == "UNDELIVERED" {
			undelivered[f[3]] = append(undelivered[f[3]], l)
		}
	}
	customer := map[string]string{} // each customer's ID, by name
	for _, rc := range simReceivers {
		customer[rc.Name] = rc.ID
	}
	idAt := map[string]time.Time{} // each event's landing
	for _, name := range order {
		sent := res.world.files[name]
		if sent == nil || sent.booking.ID == "" {
			continue
		}
		atts := byFile[name]
		if len(atts) == 0 {
			t.Errorf("no event for %s, which landed at %s, want one sent as it landed", name, landedAt[name].Format("Mon 15:04:05"))
			return
		}
		a := atts[0]
		want := Event{ID: a.Event.ID, Type: "data.landed", BookingID: sent.booking.ID, SatID: sent.booking.SatID,
			Station: sent.booking.Station, File: name, Bytes: int64(len(sent.data)), LandedAt: landedAt[name].Truncate(time.Second)}
		var took []simAttempt
		for _, b := range atts {
			if testAnswered(b) {
				took = append(took, b)
			}
		}
		sat, _ := simSatellite(sent.booking.SatID)
		last := atts[len(atts)-1].At
		switch {
		case !a.At.Equal(landedAt[name]):
			t.Errorf("%s for %s first sent at %s, want %s, as it landed", a.Event.ID, name, testAt(a.At, landedAt[name]), testAt(landedAt[name], a.At))
		case testEventDiff(a.Body, want) != "":
			t.Errorf("%s for %s: %s", a.Event.ID, name, testEventDiff(a.Body, want))
		case a.Customer != customer[sat.Operator]:
			t.Errorf("%s for %s went to %q, want %q, the customer whose name is %s's operator", a.Event.ID, name, a.Customer, customer[sat.Operator], sat.ID)
		case slices.ContainsFunc(atts, func(b simAttempt) bool { return b.Event.ID != a.Event.ID }):
			t.Errorf("%s landed once and went out as two events, %s and another, want one: an event's ID never changes", name, a.Event.ID)
		case len(took) > 0 && len(delivered[name]) != 1:
			t.Errorf("console has %d DELIVERED lines for %s, want one when its receiver took it", len(delivered[name]), name)
		case len(took) > 0 && !delivered[name][0].At.Equal(took[0].At):
			t.Errorf("DELIVERED line for %s at %s, want %s, when its receiver took it", name, testAt(delivered[name][0].At, took[0].At), testAt(took[0].At, delivered[name][0].At))
		case len(took) == 0 && len(delivered[name]) > 0:
			t.Errorf("console says DELIVERED %s, want no such line: no receiver took it", name)
		case len(took) == 0 && res.c.stopAt.IsZero() && len(undelivered[name]) != 1:
			t.Errorf("console has %d UNDELIVERED lines for %s, want one at %s: no receiver took its event", len(undelivered[name]), name, last.Format("Mon 15:04:05"))
		case len(took) == 0 && res.c.stopAt.IsZero() && !undelivered[name][0].At.Equal(last):
			t.Errorf("UNDELIVERED line for %s at %s, want %s, after its last attempt", name, testAt(undelivered[name][0].At, last), testAt(last, undelivered[name][0].At))
		default:
			if _, ok := idAt[a.Event.ID]; ok {
				t.Errorf("%s went out for two files, want an ID of its own for each", a.Event.ID)
				return
			}
			idAt[a.Event.ID] = landedAt[name]
			continue
		}
		return
	}
	for id := range idAt {
		if n, err := strconv.Atoi(strings.TrimPrefix(id, "evt-")); err != nil || n < 1 || n > len(idAt) || id != fmt.Sprintf("evt-%06d", n) {
			t.Errorf("an event is %s, want IDs from evt-000001 up to evt-%06d, one for each file", id, len(idAt))
			return
		}
	}
	for i := 1; i <= len(idAt); i++ {
		id := fmt.Sprintf("evt-%06d", i)
		at, ok := idAt[id]
		if !ok {
			t.Errorf("no event is %s, want IDs from evt-000001 up to evt-%06d, one for each file", id, len(idAt))
			return
		}
		if prev, ok := idAt[fmt.Sprintf("evt-%06d", i-1)]; ok && at.Before(prev) {
			t.Errorf("%s's file landed at %s, before evt-%06d's, want IDs in the order files land", id, testAt(at, prev), i-1)
			return
		}
	}
}

func TestDeliveries(t *testing.T) {
	res := simWeekRun(t)
	t.Run("every request signed", func(t *testing.T) {
		testRequests(t, res)
	})
	t.Run("one event per file", func(t *testing.T) {
		testEvents(t, res)
	})
	t.Run("Kessick's first", func(t *testing.T) {
		// Thistle-2's first science file lands at Port Ellery at 16:04:25 on Thursday.
		want := Event{Type: "data.landed", BookingID: "bk-8792", SatID: "thistle-2", Station: "PEL",
			File: "thistle-2-20290726-1603.dat", Bytes: 2048, LandedAt: simJuly(26, 16, 4, 25)}
		if raceOn {
			t.Skip("under -race the week flies four stations, and its bookings are numbered afresh")
		}
		for _, a := range res.world.hooks.attempts {
			if a.Event.File != want.File {
				continue
			}
			want.ID = a.Event.ID
			switch {
			case a.Customer != "kessick":
				t.Errorf("%s for %s went to %q, want Kessick's receiver", a.Event.ID, want.File, a.Customer)
			case testEventDiff(a.Body, want) != "":
				t.Errorf("%s for %s: %s", a.Event.ID, want.File, testEventDiff(a.Body, want))
			case a.Timestamp != "1879776265" || !testAnswered(a):
				t.Errorf("Kessick's receiver got %s with Apogee-Timestamp %q and answered %d, want \"1879776265\" and 204", a.Event.ID, a.Timestamp, a.Status)
			}
			return
		}
		t.Errorf("Kessick's receiver got no event for %s, want one at Thu 16:04:25, as it landed", want.File)
	})
	t.Run("a broken customers.json", func(t *testing.T) {
		c := testThursday("PEL")
		c.passes = testArchive(c, "thistle-2/PEL/02136")
		for _, tc := range []struct{ text, want string }{
			{`[{"id":"kessick","name":"Kessick Polytechnic","url":"https://hooks.kessick.example/apogee"}]` + "\n", `customers.json: json: unknown field "url"`},
			{simCustomers + "[]\n", "customers.json: more data after the customers"},
		} {
			c.customers = tc.text
			if res := simRun(t, c); res.err == nil || res.err.Error() != tc.want {
				t.Errorf("rehearse() = %v, want %q", res.err, tc.want)
				return
			} else if len(res.console) > 0 {
				t.Errorf("console = %q, want nothing after an error", res.lines())
				return
			}
		}
		c.customers = "-"
		if res := simRun(t, c); !errors.Is(res.err, fs.ErrNotExist) {
			t.Errorf("rehearse() = %v with no customers.json, want an error that wraps fs.ErrNotExist", res.err)
		}
	})
}

// testTimes writes how many times something happened: "once", "twice", "8 times".
func testTimes(n int) string {
	switch n {
	case 1:
		return "once"
	case 2:
		return "twice"
	}
	return fmt.Sprintf("%d times", n)
}

// testTries returns the requests a run's receivers got for one file's event, in order.
func testTries(res *simResult, file string) []simAttempt {
	var out []simAttempt
	for _, a := range res.world.hooks.attempts {
		if a.Event.File == file {
			out = append(out, a)
		}
	}
	return out
}

// testBackoff checks the waits between an event's attempts: after failed attempt n, from
// half of 2 s doubled for each failure before it, at most 5 minutes, to under all of it.
// It reports the first problem and returns false.
func testBackoff(t *testing.T, file string, tries []simAttempt) bool {
	t.Helper()
	for i := 1; i < len(tries); i++ {
		d := min(2*time.Second<<(i-1), 5*time.Minute)
		if gap := tries[i].At.Sub(tries[i-1].At); gap < d/2 || gap >= d {
			t.Errorf("%s for %s: attempt %d came %v after attempt %d, want from %v to under %v", tries[i].Event.ID, file, i+1, gap, i, d/2, d)
			return false
		}
	}
	return true
}

// testJitter checks that the waits between the attempts given don't all take the same
// share of their ranges: each wait's factor is drawn afresh.
func testJitter(t *testing.T, events ...[]simAttempt) {
	t.Helper()
	lo, hi, n := 2.0, 0.0, 0
	for _, tries := range events {
		for i := 1; i < len(tries); i++ {
			d := min(2*time.Second<<(i-1), 5*time.Minute)
			f := float64(tries[i].At.Sub(tries[i-1].At)) / float64(d)
			lo, hi, n = min(lo, f), max(hi, f), n+1
		}
	}
	if n > 1 && hi-lo < 1e-6 {
		t.Errorf("all %d waits between attempts took the same factor, %.3f, want one drawn at random from [0.5, 1.0) for each wait", n, lo)
	}
}

// testGivenUp checks an event its receiver never took: n attempts, with growing waits
// between them, and an UNDELIVERED line at the last. It reports the first problem only.
func testGivenUp(t *testing.T, res *simResult, file string, n int, why string) {
	t.Helper()
	tries := testTries(res, file)
	if len(tries) != n {
		t.Errorf("%s's event went out %s, want %s: %s", file, testTimes(len(tries)), testTimes(n), why)
		return
	}
	if !testBackoff(t, file, tries) {
		return
	}
	want := "UNDELIVERED " + file
	last := tries[n-1].At
	switch got := res.find(" " + want + "\n"); {
	case len(got) != 1:
		t.Errorf("console has %d lines %q, want one at %s, after the last attempt", len(got), want, last.Format("Mon 15:04:05"))
	case !got[0].At.Equal(last):
		t.Errorf("%q written at %s, want %s, after the last attempt", want, testAt(got[0].At, last), testAt(last, got[0].At))
	}
}

// testTrouble is Thursday at Port Ellery with three booked passes, each customer's receiver
// in trouble all day: Brinewatch's down, Cirrine's answering 503 and Furrowlight's 400.
func testTrouble(ids ...string) simCase {
	c := testThursday("PEL")
	c.passes, c.bookings = testArchive(c, ids...), testBooked(c, ids...)
	for _, o := range []simOutage{{Customer: "brinewatch"}, {Customer: "cirrine", Status: 503}, {Customer: "furrowlight", Status: 400}} {
		o.From, o.To = c.from, c.from.Add(48*time.Hour)
		c.receivers = append(c.receivers, o)
	}
	return c
}

func TestRetries(t *testing.T) {
	t.Run("Cirrine's receiver restarting", func(t *testing.T) {
		// Cirrine's receiver is down from 09:32:00 to 09:33:00 on Tuesday, and Cirrine-5's
		// file lands at 09:32:18.
		res := simWeekRun(t)
		const file = "cirrine-5-20290724-0931.dat"
		back := simJuly(24, 9, 33, 0)
		if tries := testTries(res, file); !raceOn && testBackoff(t, file, tries) {
			for i, a := range tries {
				switch {
				case a.At.Before(back) && testAnswered(a):
					t.Errorf("%s for %s answered at %s, want it refused until 09:33:00", a.Event.ID, file, a.At.Format("Mon 15:04:05.000"))
				case !a.At.Before(back) && !testAnswered(a):
					t.Errorf("%s for %s refused at %s, want it taken: the receiver was back at 09:33:00", a.Event.ID, file, a.At.Format("Mon 15:04:05.000"))
				case i == len(tries)-1 && !testAnswered(a):
					t.Errorf("%s for %s went out %s, the last at %s, want it sent again until the receiver takes it", a.Event.ID, file, testTimes(len(tries)), a.At.Format("Mon 15:04:05.000"))
				default:
					continue
				}
				break
			}
		}
		if got := res.find(" UNDELIVERED "); len(got) > 0 {
			t.Errorf("console says %q, want every event of the week delivered", strings.TrimSuffix(got[0].Text, "\n"))
		}
		testEmpty(t, res, "reports/undelivered.txt")
	})
	t.Run("receivers in trouble", func(t *testing.T) {
		res := simRun(t, testTrouble("brine-14/PEL/01083", "cirrine-8/PEL/01488", "furrow-a/PEL/24958"))
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testRequests(t, res)
		testGivenUp(t, res, "brine-14-20290726-0538.dat", 8, "a request that fails is sent again, 8 times in all")
		testGivenUp(t, res, "cirrine-8-20290726-0132.dat", 8, "a 503 is sent again, 8 times in all")
		testGivenUp(t, res, "furrow-a-20290726-1044.tif", 1, "a 400 is never sent again")
		if !t.Failed() {
			testJitter(t, testTries(res, "brine-14-20290726-0538.dat"), testTries(res, "cirrine-8-20290726-0132.dat"))
		}
		want := []string{
			"brine-14-20290726-0538.dat brinewatch error",
			"cirrine-8-20290726-0132.dat cirrine 503",
			"furrow-a-20290726-1044.tif furrowlight 400",
		}
		if lines, ok := res.fileLines("reports/undelivered.txt"); !ok {
			t.Errorf("reports/undelivered.txt is not there, want %q", want)
		} else if !slices.Equal(lines, want) {
			t.Errorf("reports/undelivered.txt = %q, want %q", lines, want)
		}
	})
	t.Run("the week waits for its deliveries", func(t *testing.T) {
		// Thursday at Lantern Atoll with one short pass, furrow-t/LAN/00431 (02:50:16 to
		// 02:53:35), booked for the test: its file lands at 02:51:35, and with Furrowlight's
		// receiver down the 8th attempt comes after the LOS.
		c := testThursday("LAN")
		c.passes = testArchive(c, "furrow-t/LAN/00431")
		c.bookings = []Booking{{ID: "bk-9909", SatID: "furrow-t", Station: "LAN", AOS: simJuly(26, 2, 50, 16), LOS: simJuly(26, 2, 53, 35)}}
		c.receivers = []simOutage{{Customer: "furrowlight", From: c.from, To: c.from.Add(48 * time.Hour)}}
		res := simRun(t, c)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		const file = "furrow-t-20290726-0250.tif"
		n := len(res.console)
		if n == 0 {
			t.Fatalf("console is empty, want it to end in \"week done: 1 flown, 0 failed, 0 skipped\"")
		}
		end := res.console[n-1]
		if gaveUp := res.find(" UNDELIVERED " + file + "\n"); len(gaveUp) == 0 || gaveUp[0].At.After(end.At) {
			t.Errorf("%q written at %s, before %s's event was given up: the week ends when every delivery has", strings.TrimSuffix(end.Text, "\n"), end.At.Format("Mon 15:04:05"), file)
			return
		}
		testGivenUp(t, res, file, 8, "a request that fails is sent again, 8 times in all")
		if t.Failed() {
			return
		}
		last := testTries(res, file)[7].At
		want := "week done: 1 flown, 0 failed, 0 skipped"
		switch {
		case end.Text != want+"\n":
			t.Errorf("last console line = %q, want %q", strings.TrimSuffix(end.Text, "\n"), want)
		case !end.At.Equal(last) || !res.ended.Equal(last):
			t.Errorf("%q written at %s, want %s, after the event's last attempt: the week ends when every delivery has", want, testAt(end.At, last), testAt(last, end.At))
		}
	})
	t.Run("stopped while waiting", func(t *testing.T) {
		// Cirrine-8's file lands at 01:33:35; at 01:33:55 its event waits to go again.
		c := testTrouble("cirrine-8/PEL/01488")
		c.stopAt = simJuly(26, 1, 33, 55)
		res := simRun(t, c)
		if got := res.find(" UNDELIVERED "); len(got) > 0 && got[0].At.Before(c.stopAt) {
			t.Errorf("console says %q, want no such line: a 503 is sent again, and the stop comes first", strings.TrimSuffix(got[0].Text, "\n"))
			return
		} else if len(got) > 0 {
			t.Errorf("console says %q, want no such line: a stop is not a failure", strings.TrimSuffix(got[0].Text, "\n"))
		}
		want := "rehearsal stopped at Thu 01:33:55"
		switch n := len(res.console); {
		case !res.ended.Equal(c.stopAt):
			t.Errorf("rehearse() returned at %s, want %s: a stop ends the wait before the next attempt", testAt(res.ended, c.stopAt), testAt(c.stopAt, res.ended))
		case n == 0:
			t.Errorf("console is empty, want its last line %q", want)
		case res.lines()[n-1] != want:
			t.Errorf("last console line = %q, want %q", res.lines()[n-1], want)
		}
	})
}

// testDumps returns every DUMP_FILE a run's satellites were given, in the order they got
// them.
func testDumps(res *simResult) []simCommand {
	var out []simCommand
	for _, s := range res.world.sats {
		for _, cmd := range s.commands {
			if cmd.Op == "DUMP_FILE" {
				out = append(out, cmd)
			}
		}
	}
	slices.SortFunc(out, func(a, b simCommand) int { return a.At.Compare(b.At) })
	return out
}

// testAsked checks the DUMP_FILEs of a run: exactly one, from the contact given, asking
// for the chunks given (none: the whole file). It reports the first problem only.
func testAsked(t *testing.T, res *simResult, from Contact, file string, chunks ...string) {
	t.Helper()
	want := strings.Join(append([]string{"DUMP_FILE", file}, chunks...), " ")
	landedAt := map[string]time.Time{}
	for _, l := range res.console {
		if f := strings.Fields(l.Text); len(f) == 4 && f[2] == "LANDED" {
			landedAt[f[3]] = l.At
		}
	}
	found := false
	for _, cmd := range testDumps(res) {
		got := strings.Join(append([]string{"DUMP_FILE"}, cmd.Args...), " ")
		name := ""
		if len(cmd.Args) > 0 {
			name = cmd.Args[0]
		}
		at, landed := landedAt[name]
		switch {
		case cmd.PassID == from.ID && got == want && !found:
			found = true
			continue
		case cmd.PassID == from.ID && name == file:
			t.Errorf("%s was given %q from %s, want %q", from.SatID, got, from.ID, want)
		case landed && !at.After(cmd.At):
			t.Errorf("%s was given %q from %s at %s, after the file landed at %s: want a file asked for only while incomplete", cmd.PassID[:strings.Index(cmd.PassID, "/")], got, cmd.PassID, cmd.At.Format("Mon 15:04:05"), at.Format("Mon 15:04:05"))
		case res.world.files[name] != nil:
			t.Errorf("%s was given %q from %s at %s, want none: %s was still coming down on %s", cmd.PassID[:strings.Index(cmd.PassID, "/")], got, cmd.PassID, cmd.At.Format("Mon 15:04:05"), name, res.world.files[name].passID)
		default:
			t.Errorf("%s was given %q from %s, want only %q from %s", cmd.PassID[:strings.Index(cmd.PassID, "/")], got, cmd.PassID, want, from.ID)
		}
		return
	}
	if !found {
		t.Errorf("%s was given no DUMP_FILE from %s, want %q", from.SatID, from.ID, want)
	}
}

// testNext returns a satellite's first contact in a run's plan after the one given.
func testNext(res *simResult, id string) Contact {
	plan := res.plan()
	first := plan.byID[id]
	var next Contact
	for _, c := range plan.contacts {
		if c.SatID == first.SatID && c.AOS.After(first.AOS) && (next.ID == "" || c.AOS.Before(next.AOS)) {
			next = c
		}
	}
	return next
}

func TestAskedAgain(t *testing.T) {
	t.Run("Furrow-K's image", func(t *testing.T) {
		// Furrow-K's 07:40 image over Vorrel Island on Tuesday comes down without chunks 21
		// to 23, and its next contact, over Varsk Island at 08:12:10, asks for them.
		res := simWeekRun(t)
		const file = "furrow-k-20290724-0740.tif"
		testAsked(t, res, testNext(res, "furrow-k/VRL/00403"), file, "21", "22", "23")
		if t.Failed() {
			return
		}
		if !raceOn {
			testLandedAt(t, res, file, simJuly(24, 8, 12, 30))
		}
		if tries := testTries(res, file); len(tries) != 1 || !testAnswered(tries[0]) {
			t.Errorf("%s's event went out %s, want once, taken by Furrowlight's receiver", file, testTimes(len(tries)))
		}
	})
	t.Run("a wrong CRC, asked for whole", func(t *testing.T) {
		// Thistle-2's first science file ends with a wrong CRC at Port Ellery; its next
		// contact, over Lantern Atoll at 20:12:03, asks for the whole file again.
		c := testThursday("PEL", "LAN")
		ids := []string{"thistle-2/PEL/02136", "thistle-2/LAN/02138"}
		c.passes, c.bookings = testArchive(c, ids...), testBooked(c, ids...)
		c.badCRC = []string{"thistle-2-20290726-1603.dat"}
		res := simRun(t, c)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testAsked(t, res, res.plan().byID["thistle-2/LAN/02138"], "thistle-2-20290726-1603.dat")
		testFiles(t, res)
		testEvents(t, res)
	})
	t.Run("a line dead mid-file", func(t *testing.T) {
		// Port Ellery's line dies at 16:04:20.5, after the fourth chunk of Thistle-2's first
		// science file, and stays dead past its LOS; over Lantern Atoll at 20:12:03 its next
		// contact asks for the other four.
		c := testThursday("PEL", "LAN")
		ids := []string{"thistle-2/PEL/02136", "thistle-2/LAN/02138"}
		c.passes, c.bookings = testArchive(c, ids...), testBooked(c, ids...)
		c.quiet = []simSpan{{Station: "PEL", From: simJuly(26, 16, 4, 20).Add(500 * time.Millisecond), To: simJuly(26, 16, 10, 0)}}
		res := simRun(t, c)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testAsked(t, res, res.plan().byID["thistle-2/LAN/02138"], "thistle-2-20290726-1603.dat", "4", "5", "6", "7")
		testFiles(t, res)
	})
	t.Run("a DUMP_FILE unanswered at LOS", func(t *testing.T) {
		// Lantern Atoll's line dies at 20:12:16.5, between the DUMP_FILE that asks again for
		// Thistle-2's file and its run, and stays dead past the LOS: the chunks it brings are
		// lost. On Friday at 08:02:49 the next contact finds it ran, and asks again.
		c := simCase{from: simJuly(26, 0, 0, 0), days: 2, stations: []string{"PEL", "LAN"}}
		ids := []string{"thistle-2/PEL/02136", "thistle-2/LAN/02138", "thistle-2/LAN/02146"}
		c.passes, c.bookings = testArchive(c, ids...), testBooked(c, ids...)
		c.badCRC = []string{"thistle-2-20290726-1603.dat"}
		c.quiet = []simSpan{{Station: "LAN", From: simJuly(26, 20, 12, 16).Add(500 * time.Millisecond), To: simJuly(26, 20, 19, 0)}}
		res := simRun(t, c)
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		var asked []string
		for _, cmd := range testDumps(res) {
			asked = append(asked, cmd.PassID+" "+strings.Join(cmd.Args, " "))
		}
		if want := []string{"thistle-2/LAN/02138 thistle-2-20290726-1603.dat", "thistle-2/LAN/02146 thistle-2-20290726-1603.dat"}; len(asked) == 0 {
			t.Errorf("thistle-2 was given no DUMP_FILE, want %q", want)
		} else if !slices.Equal(asked, want) {
			t.Errorf("thistle-2's DUMP_FILEs = %q, want %q: one that ran unanswered is settled on the next contact, and the file asked for again while it is incomplete", asked, want)
		}
		testFiles(t, res)
	})
	t.Run("Kessick's evt-000783", func(t *testing.T) {
		// With Furrow-K's image whole on Tuesday, Thistle-2's first science file is the
		// 783rd to land in the week.
		if raceOn {
			t.Skip("under -race the week flies four stations, and its bookings are numbered afresh")
		}
		want := Event{ID: "evt-000783", Type: "data.landed", BookingID: "bk-8792", SatID: "thistle-2", Station: "PEL",
			File: "thistle-2-20290726-1603.dat", Bytes: 2048, LandedAt: simJuly(26, 16, 4, 25)}
		const sig = "340bfc1649ec9f7fae412dddb919258dad72fc043e8bf0f45b2c8b06fd03c62e"
		tries := testTries(simWeekRun(t), want.File)
		switch {
		case len(tries) == 0:
			t.Errorf("Kessick's receiver got no event for %s, want evt-000783", want.File)
		case testEventDiff(tries[0].Body, want) != "":
			t.Errorf("Kessick's event: %s; with Furrow-K's Tuesday image landed, Thistle-2's file is the week's 783rd", testEventDiff(tries[0].Body, want))
		case tries[0].Timestamp != "1879776265" || tries[0].Signature != sig:
			t.Errorf("Kessick's event signed %s at %q, want %s at \"1879776265\"", testShort(tries[0].Signature), tries[0].Timestamp, testShort(sig))
		}
	})
}

// testRecordKey names a record of a log: its message and what it is about, the event, file
// or pass it names.
type testRecordKey struct{ msg, about string }

// testAnyText and testAnyNumber stand, in a record a test wants, for an attribute of any
// text but "" and of any number.
const (
	testAnyText   = "\x00any text"
	testAnyNumber = "\x00any number"
)

// testAttr is one attribute a test wants in a log record: a string, a number, testAnyText,
// testAnyNumber or a group of attributes.
type testAttr struct {
	key  string
	want any // string, float64 or []testAttr
}

// testRecords reads a JSON log of a run, one record a line, and returns its records by
// message and what they are about. It fails the test when the log is not there or a line
// is not a JSON object.
func testRecords(t *testing.T, res *simResult, name string) map[testRecordKey][]map[string]any {
	t.Helper()
	text, ok := res.file(name)
	if !ok {
		t.Fatalf("%s is not there, want it saved when the rehearsal ends, stopped or not", name)
	}
	out := map[testRecordKey][]map[string]any{}
	for i, line := range strings.SplitAfter(text, "\n") {
		if line == "" {
			break
		}
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("%s line %d is not one JSON object: %v", name, i+1, err)
		}
		key := testRecordKey{msg: fmt.Sprint(rec["msg"])}
		for _, k := range []string{"event_id", "file", "pass_id"} {
			if s, ok := rec[k].(string); ok {
				key.about = s
				break
			}
		}
		out[key] = append(out[key], rec)
	}
	return out
}

// testRecord returns the attributes a record wants: its time, level, message, program and
// site, then the ones given.
func testRecord(at time.Time, level, msg, program, site string, attrs ...testAttr) []testAttr {
	return append([]testAttr{{"time", at.Format(time.RFC3339Nano)}, {"level", level}, {"msg", msg},
		{"program", program}, {"site", site}}, attrs...)
}

// testPassGroup is the pass group of a contact's records.
func testPassGroup(c Contact) testAttr {
	return testAttr{"pass", []testAttr{{"sat_id", c.SatID}, {"station", c.Station},
		{"aos", c.AOS.Format(time.RFC3339Nano)}, {"los", c.LOS.Format(time.RFC3339Nano)}}}
}

// testRecordDiff names the first attribute where a record differs from want, as
// "frames = 5, want 6", or returns "".
func testRecordDiff(got map[string]any, want []testAttr, prefix string) string {
	for _, a := range want {
		v, ok := got[a.key]
		switch w := a.want.(type) {
		case []testAttr:
			g, isGroup := v.(map[string]any)
			if !isGroup {
				return fmt.Sprintf("%s%s = %s, want a group", prefix, a.key, testJSON(v))
			}
			if d := testRecordDiff(g, w, prefix+a.key+"."); d != "" {
				return d
			}
		case string:
			s, isText := v.(string)
			switch {
			case w == testAnyText && (!isText || s == ""):
				return fmt.Sprintf("%s%s = %s, want what went wrong", prefix, a.key, testJSON(v))
			case w == testAnyNumber:
				if _, isNumber := v.(float64); !isNumber {
					return fmt.Sprintf("%s%s = %s, want a number", prefix, a.key, testJSON(v))
				}
			case w != testAnyText && (!ok || v != w):
				return fmt.Sprintf("%s%s = %s, want %s", prefix, a.key, testJSON(v), testJSON(w))
			}
		default:
			if !ok || v != w {
				return fmt.Sprintf("%s%s = %s, want %s", prefix, a.key, testJSON(v), testJSON(w))
			}
		}
	}
	var extra []string
	for k := range got {
		if !slices.ContainsFunc(want, func(a testAttr) bool { return a.key == k }) {
			extra = append(extra, k)
		}
	}
	if len(extra) > 0 {
		sort.Strings(extra)
		return fmt.Sprintf("%s%s = %s, want no such attribute", prefix, extra[0], testJSON(got[extra[0]]))
	}
	return ""
}

// testWantRecords checks a log's records against the ones a test wants, in time order:
// exactly one record for each message and what it is about, with the attributes wanted,
// and no other. It reports the first problem only.
func testWantRecords(t *testing.T, name string, got map[testRecordKey][]map[string]any, want map[testRecordKey][]testAttr) {
	t.Helper()
	keys := slices.SortedFunc(maps.Keys(want), func(a, b testRecordKey) int {
		return cmp.Or(strings.Compare(want[a][0].want.(string), want[b][0].want.(string)),
			strings.Compare(a.about, b.about), strings.Compare(a.msg, b.msg))
	})
	for _, key := range keys {
		recs := got[key]
		at, _ := time.Parse(time.RFC3339Nano, want[key][0].want.(string))
		switch {
		case len(recs) == 0:
			t.Errorf("%s has no %q record for %s, want one at %s", name, key.msg, key.about, at.Format("Mon 15:04:05.999"))
		case len(recs) > 1:
			t.Errorf("%s has %d %q records for %s, want one", name, len(recs), key.msg, key.about)
		case testRecordDiff(recs[0], want[key], "") != "":
			t.Errorf("%q record for %s: %s", key.msg, key.about, testRecordDiff(recs[0], want[key], ""))
		default:
			continue
		}
		return
	}
	for _, key := range slices.SortedFunc(maps.Keys(got), func(a, b testRecordKey) int {
		return cmp.Or(strings.Compare(a.about, b.about), strings.Compare(a.msg, b.msg))
	}) {
		if _, ok := want[key]; !ok {
			t.Errorf("%s has a %q record for %s, want none", name, key.msg, key.about)
			return
		}
	}
}

// testPassRecords returns the records uplinkd's log wants for a run's contacts: "pass
// started" at each AOS the run reached, and at each LOS the record of how the contact ended,
// as its LOS line says, unless the run was stopped first. A pass flown counts its TLM frames
// with a good CRC, any number of them at a station that flips bytes, and the commands it
// sent that its satellite ran.
func testPassRecords(res *simResult) map[testRecordKey][]testAttr {
	ran := map[string]float64{} // by pass ID
	for _, s := range res.world.sats {
		for _, cmd := range s.commands {
			if cmd.Refused == "" {
				ran[cmd.PassID]++
			}
		}
	}
	outcomes := res.outcomes()
	want := map[testRecordKey][]testAttr{}
	for _, c := range res.plan().contacts {
		o, ok := outcomes[c.ID]
		if !ok {
			continue
		}
		group, id := testPassGroup(c), testAttr{"pass_id", c.ID}
		want[testRecordKey{"pass started", c.ID}] = testRecord(c.AOS, "INFO", "pass started", "uplinkd", c.Station,
			group, id, testAttr{"dish", float64(c.Dish)})
		kind, reason, _ := strings.Cut(o, " ")
		switch kind {
		case "flown":
			var frames any = testAnyNumber
			if p := res.world.passes[c.ID]; p != nil && res.c.flipEvery[c.Station] == 0 {
				frames = float64(p.tlm)
			}
			want[testRecordKey{"pass flown", c.ID}] = testRecord(c.LOS, "INFO", "pass flown", "uplinkd", c.Station,
				group, id, testAttr{"max_elev_deg", c.MaxElevDeg}, testAttr{"frames", frames}, testAttr{"commands", ran[c.ID]})
		case "failed":
			want[testRecordKey{"pass failed", c.ID}] = testRecord(c.LOS, "ERROR", "pass failed", "uplinkd", c.Station,
				group, id, testAttr{"max_elev_deg", c.MaxElevDeg}, testAttr{"reason", reason}, testAttr{"err", testAnyText})
		case "skipped":
			want[testRecordKey{"pass skipped", c.ID}] = testRecord(c.LOS, "WARN", "pass skipped", "uplinkd", c.Station,
				group, id, testAttr{"reason", reason})
		}
	}
	return want
}

// testFileRecords returns the records dlinkd's log wants for a run's console: "file
// landed" at each LANDED line, on the pass whose call brought its last chunk, "event
// delivered" at each DELIVERED line and "event undelivered" at each UNDELIVERED line, with
// the attempts the event's receiver saw.
func testFileRecords(res *simResult) map[testRecordKey][]testAttr {
	landedOn := map[string]string{} // by file
	for _, ch := range res.world.chunks {
		landedOn[ch.File] = ch.PassID
	}
	tries := map[string][]simAttempt{} // by file
	for _, a := range res.world.hooks.attempts {
		tries[a.Event.File] = append(tries[a.Event.File], a)
	}
	want := map[testRecordKey][]testAttr{}
	for _, l := range res.console {
		f := strings.Fields(l.Text)
		if len(f) != 4 {
			continue
		}
		name := f[3]
		sent, atts := res.world.files[name], tries[name]
		switch {
		case f[2] == "LANDED" && sent != nil:
			want[testRecordKey{"file landed", name}] = testRecord(l.At, "INFO", "file landed", "dlinkd", "HQ",
				testAttr{"file", name}, testAttr{"bytes", float64(len(sent.data))}, testAttr{"booking_id", sent.booking.ID},
				testAttr{"pass_id", landedOn[name]})
		case f[2] == "DELIVERED" && len(atts) > 0:
			id := atts[0].Event.ID
			want[testRecordKey{"event delivered", id}] = testRecord(l.At, "INFO", "event delivered", "dlinkd", "HQ",
				testAttr{"event_id", id}, testAttr{"customer", atts[0].Customer}, testAttr{"attempts", float64(len(atts))})
		case f[2] == "UNDELIVERED" && len(atts) > 0:
			id := atts[0].Event.ID
			want[testRecordKey{"event undelivered", id}] = testRecord(l.At, "ERROR", "event undelivered", "dlinkd", "HQ",
				testAttr{"event_id", id}, testAttr{"customer", atts[0].Customer}, testAttr{"attempts", float64(len(atts))},
				testAttr{"err", testAnyText})
		}
	}
	return want
}

// testSecrets checks that no log of a run holds a secret: a master's text, hex or base64, a
// satellite's key derived from one, in hex or base64, or a webhook secret. It reports the
// first one found.
func testSecrets(t *testing.T, res *simResult) {
	t.Helper()
	texts := map[string]string{} // each secret's text, and what it is
	keys := map[string]string{}  // each key's hex and base64, and what it is
	for _, m := range simParseKeys(simKeys) {
		texts[string(m.Secret)] = "master " + m.ID
		keys[hex.EncodeToString(m.Secret)] = "master " + m.ID + " in hex"
		keys[base64.StdEncoding.EncodeToString(m.Secret)] = "master " + m.ID + " in base64"
		for _, s := range simNet().Satellites {
			for _, purpose := range []string{"sign", "seal"} {
				k := simSatKey(m.Secret, purpose, s.ID)
				what := fmt.Sprintf("%s's %s key from %s", s.ID, purpose, m.ID)
				keys[hex.EncodeToString(k)] = what + " in hex"
				keys[base64.StdEncoding.EncodeToString(k)] = what + " in base64"
			}
		}
	}
	for _, rc := range simReceivers {
		texts[rc.Secret] = rc.ID + "'s webhook secret"
	}
	if len(res.names("logs")) == 0 {
		t.Errorf("logs/ holds no log, want the rehearsal's logs there")
		return
	}
	isKeyChar := func(b byte) bool {
		return '0' <= b && b <= '9' || 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '+' || b == '/' || b == '='
	}
	for _, name := range res.names("logs") {
		text, _ := res.file("logs/" + name)
		for _, secret := range slices.Sorted(maps.Keys(texts)) {
			if strings.Contains(text, secret) {
				t.Errorf("logs/%s holds %s, %q, want no secret in any log", name, texts[secret], secret)
				return
			}
		}
		// A key hides in a run of hex or base64 digits at least 43 long; a record has none
		// of its own.
		start := 0
		for i := 0; i <= len(text); i++ {
			if i < len(text) && isKeyChar(text[i]) {
				continue
			}
			if run := text[start:i]; len(run) >= 43 {
				for _, k := range slices.Sorted(maps.Keys(keys)) {
					if strings.Contains(run, k) || strings.Contains(strings.ToLower(run), k) {
						t.Errorf("logs/%s holds %s, want no key in any log", name, keys[k])
						return
					}
				}
			}
			start = i + 1
		}
	}
}

func TestLogs(t *testing.T) {
	res := simWeekRun(t)
	if res.err != nil {
		t.Fatalf("rehearse() = %v, want nil", res.err)
	}
	t.Run("uplinkd's week", func(t *testing.T) {
		got := testRecords(t, res, "logs/uplinkd.log")
		want := testPassRecords(res)
		// Ystmark's fiber dies at 03:31 on Thursday; the last frame of pipit-2/YST/22861
		// came at 03:30:52, so its line is lost at 03:32:22.
		want[testRecordKey{"line lost", "pipit-2/YST/22861"}] = testRecord(simJuly(26, 3, 32, 22), "WARN", "line lost", "uplinkd", "YST",
			testAttr{"pass_id", "pipit-2/YST/22861"}, testAttr{"err", testAnyText})
		testWantRecords(t, "logs/uplinkd.log", got, want)
		if t.Failed() {
			return
		}
		// Varsk Island flips a byte in one TLM frame in 400, which the frames of its passes
		// leave out.
		sent, counted := 0, 0
		for key, recs := range got {
			if key.msg == "pass flown" && res.c.flipEvery[recs[0]["site"].(string)] > 0 {
				sent += res.world.passes[key.about].tlm
				counted += int(recs[0]["frames"].(float64))
			}
		}
		if counted != sent-res.world.tlmBad {
			t.Errorf("the passes flown at Varsk Island counted %d TLM frames, want %d: %d sent, %d with a flipped byte", counted, sent-res.world.tlmBad, sent, res.world.tlmBad)
		}
	})
	t.Run("dlinkd's week", func(t *testing.T) {
		testWantRecords(t, "logs/dlinkd.log", testRecords(t, res, "logs/dlinkd.log"), testFileRecords(res))
	})
	t.Run("nothing secret", func(t *testing.T) {
		testSecrets(t, res)
	})
	t.Run("receivers in trouble", func(t *testing.T) {
		res := simRun(t, testTrouble("brine-14/PEL/01083", "cirrine-8/PEL/01488", "furrow-a/PEL/24958"))
		if res.err != nil {
			t.Fatalf("rehearse() = %v, want nil", res.err)
		}
		testWantRecords(t, "logs/dlinkd.log", testRecords(t, res, "logs/dlinkd.log"), testFileRecords(res))
	})
	t.Run("stopped at noon", func(t *testing.T) {
		// At Tuesday 12:00:00 furrow-f/YST/00406 is flying: it has no end record.
		stop := simJuly(24, 12, 0, 0)
		res := simRun(t, simCase{from: simJuly(24, 0, 0, 0), days: 1, stations: []string{"YST"}, stopAt: stop})
		testWantRecords(t, "logs/uplinkd.log", testRecords(t, res, "logs/uplinkd.log"), testPassRecords(res))
	})
}
