package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/iotest"
	"time"
)

// session runs s on input and returns what it wrote and what Run returned.
func session(s *Shell, input string) (string, error) {
	var out strings.Builder
	err := s.Run(strings.NewReader(input), &out)
	return out.String(), err
}

// testCommand returns a test-only command that writes its arguments.
func testCommand(help string) command {
	return command{help: help, run: func(s *Shell, args []string, out io.Writer) error {
		fmt.Fprintf(out, "ran with %q\n", args)
		return nil
	}}
}

// night returns h:m on the night of the passes, UTC.
func night(h, m int) time.Time { return time.Date(2027, time.October, 30, h, m, 0, 0, time.UTC) }

// fixed returns a clock that always says at.
func fixed(at time.Time) func() time.Time { return func() time.Time { return at } }

// check runs input on a new shell that also knows a test-only wind command, and compares what
// it wrote before the closing bye with want; TestExit checks the bye.
func check(t *testing.T, input, want string) {
	t.Helper()
	s := newShell()
	s.cache.now = fixed(night(4, 0))
	s.commands["wind"] = testCommand("read Cape Wren's wind")
	got, err := session(s, input)
	got = strings.TrimSuffix(got, "bye\n")
	if got != want {
		t.Errorf("Run(%q) wrote %q, want %q", input, got, want)
	}
	if err != nil {
		t.Errorf("Run(%q) = %v, want nil", input, err)
	}
}

func TestPrompt(t *testing.T) {
	for _, tc := range []struct{ name, input, want string }{
		{"one line", "track\n", "gctl> unknown command \"track\"; type help\ngctl> "},
		{"two lines", "track\npoint\n", "gctl> unknown command \"track\"; type help\ngctl> unknown command \"point\"; type help\ngctl> "},
		{"no newline at the end", "track", "gctl> unknown command \"track\"; type help\ngctl> "},
		{"nothing typed", "", "gctl> "},
	} {
		t.Run(tc.name, func(t *testing.T) { check(t, tc.input, tc.want) })
	}
}

func TestWords(t *testing.T) {
	unknown := "gctl> unknown command \"track\"; type help\ngctl> "
	for _, tc := range []struct{ name, input, want string }{
		{"two spaces", "track  pipit-1\n", unknown},
		{"tabs and spaces", "\t track \tpipit-1   PEL \n", unknown},
		{"blank line", "\n", "gctl> gctl> "},
		{"spaces and tabs only", "  \t \n", "gctl> gctl> "},
		{"blank line between", "track\n\npoint\n", "gctl> unknown command \"track\"; type help\ngctl> gctl> unknown command \"point\"; type help\ngctl> "},
	} {
		t.Run(tc.name, func(t *testing.T) { check(t, tc.input, tc.want) })
	}
}

func TestHelp(t *testing.T) {
	check(t, "help\n", "gctl> "+
		"exit      leave the shell\n"+
		"help      list every command\n"+
		"history   list the commands typed so far\n"+
		"next      show the next pass\n"+
		"passes    list tonight's passes\n"+
		"stations  list the stations and which are stowed\n"+
		"stow      record a stow: stow CODE [reason...]\n"+
		"wind      read Cape Wren's wind\n"+
		"gctl> ")

	s := newShell()
	for _, c := range []struct{ name, help string }{
		{"wind", "read Cape Wren's wind"},
		{"uplink", "send a signed command"},
		{"doppler", "show a beacon's frequency shift"},
		{"beacon", "listen for a satellite's beacon"},
		{"azimuth", "turn the dish to a compass direction"},
	} {
		s.commands[c.name] = testCommand(c.help)
	}
	got, _ := session(s, "help please\n")
	want := "gctl> " +
		"azimuth   turn the dish to a compass direction\n" +
		"beacon    listen for a satellite's beacon\n" +
		"doppler   show a beacon's frequency shift\n" +
		"exit      leave the shell\n" +
		"help      list every command\n" +
		"history   list the commands typed so far\n" +
		"next      show the next pass\n" +
		"passes    list tonight's passes\n" +
		"stations  list the stations and which are stowed\n" +
		"stow      record a stow: stow CODE [reason...]\n" +
		"uplink    send a signed command\n" +
		"wind      read Cape Wren's wind\n" +
		"gctl> "
	if strings.TrimSuffix(got, "bye\n") != want {
		t.Errorf("Run(%q) with five test commands wrote %q, want %q", "help please\n", got, want)
	}
}

func TestCommands(t *testing.T) {
	for _, tc := range []struct{ name, input, want string }{
		{"arguments", "wind CWR  gusts\n", "gctl> ran with [\"CWR\" \"gusts\"]\ngctl> "},
		{"no arguments", "wind\n", "gctl> ran with []\ngctl> "},
		{"unknown command", "track CWR\n", "gctl> unknown command \"track\"; type help\ngctl> "},
	} {
		t.Run(tc.name, func(t *testing.T) { check(t, tc.input, tc.want) })
	}
}

func TestExit(t *testing.T) {
	for _, tc := range []struct{ name, input, want string }{
		{"exit", "exit\n", "gctl> bye\n"},
		{"exit with words after it", "exit now\n", "gctl> bye\n"},
		{"lines after exit", "exit\nwind CWR\n", "gctl> bye\n"},
		{"end of input", "wind\n", "gctl> ran with []\ngctl> bye\n"},
		{"nothing typed", "", "gctl> bye\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := newShell()
			s.commands["wind"] = testCommand("read Cape Wren's wind")
			got, err := session(s, tc.input)
			if got != tc.want {
				t.Errorf("Run(%q) wrote %q, want %q", tc.input, got, tc.want)
			}
			if err != nil {
				t.Errorf("Run(%q) = %v, want nil", tc.input, err)
			}
		})
	}
}

func TestCommandError(t *testing.T) {
	s := newShell()
	s.commands["point"] = command{help: "point the dish", run: func(s *Shell, args []string, out io.Writer) error {
		return errors.New("dish stowed")
	}}
	input := "point CWR\npoint\n"
	got, err := session(s, input)
	want := "gctl> error: dish stowed\ngctl> error: dish stowed\ngctl> bye\n"
	if got != want {
		t.Errorf("Run(%q) wrote %q, want %q", input, got, want)
	}
	if err != nil {
		t.Errorf("Run(%q) = %v, want nil: a command's error is not the session's", input, err)
	}
}

func TestReadError(t *testing.T) {
	problem := errors.New("terminal closed")
	var out strings.Builder
	err := newShell().Run(io.MultiReader(strings.NewReader("track\n"), iotest.ErrReader(problem)), &out)
	if !errors.Is(err, problem) {
		t.Errorf("Run(input failing after one line) = %v, want %v", err, problem)
	}
	want := "gctl> unknown command \"track\"; type help\ngctl> "
	if out.String() != want {
		t.Errorf("Run(input failing after one line) wrote %q, want %q, with no bye", out.String(), want)
	}
}

// tonightLines is what passes prints for the whole night.
const tonightLines = "04:12-04:21 PEL pipit-1 max 61.4\n" +
	"04:31-04:40 YST pipit-2 max 47.2\n" +
	"04:58-05:06 CWR furrow-a max 33.5\n" +
	"05:20-05:26 LAN brine-3 max 22.8\n" +
	"05:47-05:57 YST cirrine-1 max 71.9\n"

func TestPasses(t *testing.T) {
	check(t, "passes\n", "gctl> "+tonightLines+"gctl> ")
	check(t, "passes\npasses\n", "gctl> "+tonightLines+"gctl> "+tonightLines+"gctl> ")
}

func TestPassesAtStation(t *testing.T) {
	for _, tc := range []struct{ name, input, want string }{
		{"Ystmark", "passes YST\n", "gctl> 04:31-04:40 YST pipit-2 max 47.2\n05:47-05:57 YST cirrine-1 max 71.9\ngctl> "},
		{"Port Ellery", "passes PEL\n", "gctl> 04:12-04:21 PEL pipit-1 max 61.4\ngctl> "},
		{"no passes tonight", "passes SQT\n", "gctl> gctl> "},
		{"unknown code", "passes XYZ\n", "gctl> error: unknown station \"XYZ\"\ngctl> "},
		{"lowercase code", "passes pel\n", "gctl> error: unknown station \"pel\"\ngctl> "},
		{"two codes", "passes PEL YST\n", "gctl> error: passes takes at most one station code\ngctl> "},
		{"the session goes on", "passes XYZ\npasses PEL\n", "gctl> error: unknown station \"XYZ\"\ngctl> 04:12-04:21 PEL pipit-1 max 61.4\ngctl> "},
	} {
		t.Run(tc.name, func(t *testing.T) { check(t, tc.input, tc.want) })
	}
}

// stationLines is what stations prints when no station is stowed.
const stationLines = "YST Ystmark\n" +
	"PEL Port Ellery\n" +
	"SQT Salar Quieto\n" +
	"LAN Lantern Atoll\n" +
	"CWR Cape Wren\n"

func TestStations(t *testing.T) {
	check(t, "stations\n", "gctl> "+stationLines+"gctl> ")
	check(t, "stations YST\n", "gctl> error: stations takes no arguments\ngctl> ")

	s := newShell()
	cwr, sqt := &s.stations[stationIndex(s.stations, "CWR")], &s.stations[stationIndex(s.stations, "SQT")]
	cwr.Stowed, cwr.StowReason = true, "wind 74 km/h"
	sqt.Stowed = true
	got, _ := session(s, "stations\n")
	want := "gctl> " +
		"YST Ystmark\n" +
		"PEL Port Ellery\n" +
		"SQT Salar Quieto stowed\n" +
		"LAN Lantern Atoll\n" +
		"CWR Cape Wren stowed: wind 74 km/h\n" +
		"gctl> "
	if strings.TrimSuffix(got, "bye\n") != want {
		t.Errorf("Run(%q) with SQT and CWR stowed wrote %q, want %q", "stations\n", got, want)
	}
}

func TestStow(t *testing.T) {
	for _, tc := range []struct{ name, input, want string }{
		{"stow with a reason", "stow CWR wind 74 km/h\nstations\n", "gctl> CWR stowed: wind 74 km/h\ngctl> YST Ystmark\nPEL Port Ellery\nSQT Salar Quieto\nLAN Lantern Atoll\nCWR Cape Wren stowed: wind 74 km/h\ngctl> "},
		{"reason words joined", "stow CWR  wind\t74  km/h\n", "gctl> CWR stowed: wind 74 km/h\ngctl> "},
		{"reason as typed", "stow SQT Dust on La Grande\n", "gctl> SQT stowed: Dust on La Grande\ngctl> "},
		{"no reason", "stow PEL\nstations\n", "gctl> PEL stowed\ngctl> YST Ystmark\nPEL Port Ellery stowed\nSQT Salar Quieto\nLAN Lantern Atoll\nCWR Cape Wren\ngctl> "},
		{"no code", "stow\n", "gctl> error: stow needs a station code\ngctl> "},
		{"unknown code", "stow XYZ high wind\nstations\n", "gctl> error: unknown station \"XYZ\"\ngctl> YST Ystmark\nPEL Port Ellery\nSQT Salar Quieto\nLAN Lantern Atoll\nCWR Cape Wren\ngctl> "},
		{"already stowed", "stow CWR wind\nstow CWR gusts\nstations\n", "gctl> CWR stowed: wind\ngctl> error: CWR is already stowed\ngctl> YST Ystmark\nPEL Port Ellery\nSQT Salar Quieto\nLAN Lantern Atoll\nCWR Cape Wren stowed: wind\ngctl> "},
	} {
		t.Run(tc.name, func(t *testing.T) { check(t, tc.input, tc.want) })
	}
}

func TestPanic(t *testing.T) {
	s := newShell()
	s.commands["doppler"] = command{help: "show a beacon's frequency shift", run: func(s *Shell, args []string, out io.Writer) error {
		var shifts map[string]float64
		shifts["pipit-1"] = 1.2 // a bug: the map was never made
		return nil
	}}
	s.commands["point"] = command{help: "point the dish", run: func(s *Shell, args []string, out io.Writer) error {
		return errors.New("dish stowed")
	}}
	input := "doppler\npasses PEL\npoint\ndoppler pipit-1\nexit\n"
	want := "gctl> error: internal error in doppler\n" +
		"gctl> 04:12-04:21 PEL pipit-1 max 61.4\n" +
		"gctl> error: dish stowed\n" +
		"gctl> error: internal error in doppler\n" +
		"gctl> bye\n"
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Run(%q) panicked with %v, want the error line and the session going on", input, r)
		}
	}()
	got, err := session(s, input)
	if got != want {
		t.Errorf("Run(%q) wrote %q, want %q", input, got, want)
	}
	if err != nil {
		t.Errorf("Run(%q) = %v, want nil", input, err)
	}
}

// countingLoader returns a loader that counts its calls in calls and gives the night's passes.
func countingLoader(calls *int) func() ([]Pass, error) {
	return func() ([]Pass, error) {
		*calls++
		return tonight(), nil
	}
}

func TestLoadOnce(t *testing.T) {
	calls := 0
	c := &passCache{now: fixed(night(4, 0)), ttl: 10 * time.Minute, load: countingLoader(&calls)}
	for i := 1; i <= 3; i++ {
		ps, err := c.passes()
		if len(ps) != 5 || err != nil {
			t.Errorf("passes() call %d = %d passes, %v, want 5, nil", i, len(ps), err)
		}
	}
	if calls != 1 {
		t.Errorf("3 calls to passes() called load %d times, want 1", calls)
	}
}

func TestLoadRetry(t *testing.T) {
	problem := errors.New("schedule file unreachable")
	calls := 0
	c := &passCache{now: fixed(night(4, 0)), ttl: 10 * time.Minute, load: func() ([]Pass, error) {
		calls++
		if calls == 1 {
			return nil, problem
		}
		return tonight(), nil
	}}
	ps, err := c.passes()
	if ps != nil || !errors.Is(err, problem) {
		t.Errorf("passes() with a failing load = %v, %v, want nil, %v", ps, err, problem)
	}
	ps, err = c.passes()
	if len(ps) != 5 || err != nil {
		t.Errorf("passes() after a failed load = %d passes, %v, want 5, nil: try again", len(ps), err)
	}
	c.passes()
	if calls != 2 {
		t.Errorf("load was called %d times, want 2: one failure, then one load kept", calls)
	}
}

func TestCacheInShell(t *testing.T) {
	s := newShell()
	if s.cache.loaded {
		t.Errorf("newShell() loaded the passes, want them loaded on first use")
	}
	calls, at := 0, night(4, 0)
	s.cache.load = countingLoader(&calls)
	s.cache.now = func() time.Time { return at }
	for _, step := range []struct {
		at    time.Time
		loads int
	}{{night(4, 0), 1}, {night(4, 9), 1}, {night(4, 10), 2}} {
		at = step.at
		session(s, "passes\npasses PEL\n")
		if calls != step.loads {
			t.Errorf("Run(%q) at %s: load called %d times in all, want %d: fresh for 10 minutes", "passes\npasses PEL\n", at.Format("15:04"), calls, step.loads)
		}
	}
	s.cache.load = func() ([]Pass, error) { return nil, errors.New("schedule file unreachable") }
	at = night(4, 30)
	got, _ := session(s, "passes\n")
	if want := "gctl> error: schedule file unreachable\ngctl> bye\n"; got != want {
		t.Errorf("Run(%q) with a failing load wrote %q, want %q", "passes\n", got, want)
	}
}

func TestStale(t *testing.T) {
	calls, at := 0, night(4, 0)
	c := &passCache{load: countingLoader(&calls), now: func() time.Time { return at }, ttl: 10 * time.Minute}
	for _, step := range []struct {
		name  string
		at    time.Time
		loads int
	}{
		{"first use at 04:00", night(4, 0), 1},
		{"fresh at 04:09:59", night(4, 9).Add(59 * time.Second), 1},
		{"stale at exactly 04:10", night(4, 10), 2},
		{"fresh again at 04:19", night(4, 19), 2},
		{"stale at 04:20", night(4, 20), 3},
	} {
		at = step.at
		if _, err := c.passes(); err != nil {
			t.Fatalf("passes() %s = %v, want nil", step.name, err)
		}
		if calls != step.loads {
			t.Errorf("passes() %s: load called %d times in all, want %d", step.name, calls, step.loads)
		}
	}
}

func TestNext(t *testing.T) {
	for _, tc := range []struct {
		name string
		at   time.Time
		want string
	}{
		{"before the first pass", night(4, 0), "next: pipit-1 over PEL at 04:12, in 12m0s"},
		{"between passes", night(4, 25), "next: pipit-2 over YST at 04:31, in 6m0s"},
		{"exactly at an AOS", night(4, 12), "next: pipit-2 over YST at 04:31, in 19m0s"},
		{"a second after an AOS", night(4, 58).Add(time.Second), "next: brine-3 over LAN at 05:20, in 21m59s"},
		{"at the last AOS", night(5, 47), "no more passes tonight"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := newShell()
			s.cache.now = fixed(tc.at)
			got, err := session(s, "next\n")
			if want := "gctl> " + tc.want + "\ngctl> bye\n"; got != want || err != nil {
				t.Errorf("Run(%q) at %s wrote %q, %v, want %q, nil", "next\n", tc.at.Format("15:04:05"), got, err, want)
			}
		})
	}
	check(t, "next YST\n", "gctl> error: next takes no arguments\ngctl> ")
}

func TestRefresher(t *testing.T) {
	var running, overlaps atomic.Int32
	c := &passCache{now: fixed(night(4, 0)), ttl: 10 * time.Minute, load: func() ([]Pass, error) {
		if running.Add(1) > 1 {
			overlaps.Add(1)
		}
		time.Sleep(time.Millisecond)
		running.Add(-1)
		return tonight(), nil
	}}
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 5 {
				c.refresh()
				if ps, err := c.passes(); len(ps) != 5 || err != nil {
					t.Errorf("passes() beside a refresher = %d passes, %v, want 5, nil", len(ps), err)
				}
			}
		})
	}
	wg.Wait()
	if n := overlaps.Load(); n > 0 {
		t.Errorf("load started %d times while another load was running, want never: lock the cache", n)
	}

	ps, _ := c.passes()
	ps[0].SatID = "changed"
	if again, _ := c.passes(); again[0].SatID == "changed" {
		t.Errorf("passes() returned the cache's own slice: a caller changed the cache, want a copy")
	}
}

// writeSettings writes text to a gctl.json in a new temporary folder and returns its path.
func writeSettings(t *testing.T, text string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "gctl.json")
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadSettings(t *testing.T) {
	got, err := loadSettings(filepath.Join(t.TempDir(), "gctl.json"))
	if got != defaultSettings() || err != nil {
		t.Errorf("loadSettings(missing file) = %+v, %v, want the defaults %+v, nil", got, err, defaultSettings())
	}
	ttl := func(d time.Duration) ttlText { return ttlText(d) }
	for _, tc := range []struct {
		name, text string
		want       Settings
	}{
		{"empty object", `{}`, Settings{Prompt: "gctl> ", CacheTTL: ttl(10 * time.Minute)}},
		{"station only", `{"station":"PEL"}`, Settings{Prompt: "gctl> ", CacheTTL: ttl(10 * time.Minute), Station: "PEL"}},
		{"every field", `{"prompt":"night> ","cache_ttl":"5m","station":"PEL"}`, Settings{Prompt: "night> ", CacheTTL: ttl(5 * time.Minute), Station: "PEL"}},
		{"spaces and lines after", "{\"cache_ttl\": \"90s\"}\n\n", Settings{Prompt: "gctl> ", CacheTTL: ttl(90 * time.Second)}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := loadSettings(writeSettings(t, tc.text))
			if got != tc.want || err != nil {
				t.Errorf("loadSettings(%q) = %+v, %v, want %+v, nil", tc.text, got, err, tc.want)
			}
		})
	}
	for _, tc := range []struct{ name, text string }{
		{"misspelled field", `{"promt":"night> "}`},
		{"two objects", `{"station":"PEL"} {"station":"YST"}`},
		{"text after the object", `{"station":"PEL"} thanks`},
		{"duration as a number", `{"cache_ttl":600}`},
		{"duration without a unit", `{"cache_ttl":"10"}`},
		{"zero duration", `{"cache_ttl":"0s"}`},
		{"not JSON", "station: PEL\n"},
		{"empty file", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := loadSettings(writeSettings(t, tc.text)); err == nil {
				t.Errorf("loadSettings(%q) error = nil, want an error", tc.text)
			}
		})
	}
}

func TestApply(t *testing.T) {
	s := newShell()
	s.cache.now = fixed(night(4, 0))
	s.apply(Settings{Prompt: "night> ", CacheTTL: ttlText(5 * time.Minute), Station: "PEL"})
	got, _ := session(s, "passes\npasses YST\n")
	want := "night> 04:12-04:21 PEL pipit-1 max 61.4\n" +
		"night> 04:31-04:40 YST pipit-2 max 47.2\n05:47-05:57 YST cirrine-1 max 71.9\n" +
		"night> bye\n"
	if got != want {
		t.Errorf("Run(%q) after apply wrote %q, want %q", "passes\npasses YST\n", got, want)
	}
	if s.cache.ttl != 5*time.Minute {
		t.Errorf("after apply, s.cache.ttl = %v, want 5m0s", s.cache.ttl)
	}

	s = newShell()
	s.cache.now = fixed(night(4, 0))
	s.apply(defaultSettings())
	if got, _ := session(s, "passes\n"); got != "gctl> "+tonightLines+"gctl> bye\n" {
		t.Errorf("Run(%q) after applying the defaults wrote %q, want every pass", "passes\n", got)
	}

	s = newShell()
	s.cache.now = fixed(night(4, 0))
	s.apply(Settings{Prompt: "gctl> ", CacheTTL: ttlText(10 * time.Minute), Station: "XYZ"})
	if got, _ := session(s, "passes\n"); got != "gctl> error: unknown station \"XYZ\"\ngctl> bye\n" {
		t.Errorf("Run(%q) with the default station XYZ wrote %q, want the unknown-station error", "passes\n", got)
	}
}

func TestHistory(t *testing.T) {
	check(t, "passes PEL\n\ntrack\nhistory\n", "gctl> 04:12-04:21 PEL pipit-1 max 61.4\n"+
		"gctl> gctl> unknown command \"track\"; type help\n"+
		"gctl>   1  passes PEL\n  2  track\n  3  history\n"+
		"gctl> ")
	check(t, strings.Repeat("track\n", 10)+"history\n", "gctl> "+
		"unknown command \"track\"; type help\ngctl> unknown command \"track\"; type help\ngctl> unknown command \"track\"; type help\ngctl> unknown command \"track\"; type help\ngctl> unknown command \"track\"; type help\ngctl> unknown command \"track\"; type help\ngctl> unknown command \"track\"; type help\ngctl> unknown command \"track\"; type help\ngctl> unknown command \"track\"; type help\ngctl> unknown command \"track\"; type help\ngctl> "+
		"  1  track\n  2  track\n  3  track\n  4  track\n  5  track\n  6  track\n  7  track\n  8  track\n  9  track\n 10  track\n 11  history\n"+
		"gctl> ")
	check(t, "stow CWR  wind\nhistory 5\n", "gctl> CWR stowed: wind\n"+
		"gctl> error: history takes no arguments\n"+
		"gctl> ")
	check(t, "stow CWR  wind\nhistory\n", "gctl> CWR stowed: wind\n"+
		"gctl>   1  stow CWR  wind\n  2  history\n"+
		"gctl> ")
}
