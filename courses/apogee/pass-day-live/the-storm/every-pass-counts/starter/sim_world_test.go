package main

// The simulated world's harness. A test describes one rehearsal with a simCase and runs it
// with simRun: the harness writes the rehearsal's folder, builds a World whose clock, radios
// and receivers are simulated, runs rehearse inside a testing/synctest bubble, so a week
// passes in a moment, and returns everything that happened. Every top-level name in the
// sim files begins with sim, so none can collide with a name in the rehearsal.

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// simSpan is a stretch of time at one station: From included, To excluded. A zero From is
// the beginning of time, and a zero To never comes.
type simSpan struct {
	Station  string
	From, To time.Time
}

// has reports whether the span holds an instant at a station.
func (s simSpan) has(station string, at time.Time) bool {
	return s.Station == station && (s.From.IsZero() || !at.Before(s.From)) && (s.To.IsZero() || at.Before(s.To))
}

// simSlow is a station's slow line for a while: every message takes Delay longer each way,
// and the station's own arrive in pieces of at most Piece bytes, Gap apart.
type simSlow struct {
	Station  string
	From, To time.Time
	Delay    time.Duration
	Piece    int
	Gap      time.Duration
}

// simOutage is a customer's receiver in trouble for a while: with Status 0 every request
// to it fails, and otherwise it answers that status.
type simOutage struct {
	Customer string // "cirrine"
	From, To time.Time
	Status   int
}

// simCase is one rehearsal a test runs: what it flies, what its folder holds and the
// trouble the simulated world makes. Whatever a case leaves out is off.
type simCase struct {
	// What the rehearsal flies.
	from     time.Time // the window's first instant, 00:00 UTC of a day; zero is the investors' week
	days     int       // 1 to 7; zero is 7
	stations []string  // station codes, in any order; none is all twenty

	// The rehearsal's folder.
	settings     string    // rehearsal.json as text; "" is made from the fields above, "-" writes none
	passes       string    // passes.txt as text; "" is passplan's archive for the window, "-" writes none
	archive      []string  // the stations whose passes the archive holds, when not just the rehearsed ones
	bookings     []Booking // bookings.json; nil writes none
	bookingsText string    // bookings.json as text, instead of bookings
	keys         string    // uplink-keys.txt; "" is simKeys, the investors' week's; "-" writes none
	customers    string    // customers.json; "" is simCustomers, the five customers'; "-" writes none

	// The trouble.
	bannerOnly bool                     // stations send their banner, then nothing, and never hang up
	silent     []string                 // stations that answer a call and never send a byte
	banner     map[string]string        // a station's banner, instead of "<CODE> READY PROTO 1"
	dialFail   []simSpan                // a station's calls fail from From to To (zero times: all week)
	split      map[string]int           // stations that write every message in pieces of at most n bytes
	glue       []string                 // stations that write the banner and their first frame in one write
	corrupt    map[string]int           // pass ID to n: the pass's n-th TLM frame arrives with a flipped byte
	flipEvery  map[string]int           // station to n: every n-th TLM frame it relays arrives with a flipped byte
	badLength  map[string]int           // pass ID to n: a frame with length 8 follows the pass's n-th frame
	late       map[string]time.Duration // stations that relay each TLM frame this long after it was made
	quiet      []simSpan                // a station's line dead: calls open at From go silent for good, new ones fail
	backhaul   []simSlow                // a station's line slow
	dropBefore map[string]int           // pass ID to n: its call drops just before the satellite gets its n-th command
	dropAfter  map[string]int           // pass ID to n: its call drops just after the satellite runs its n-th command
	stows      []Stow                   // wind locks: a station's dishes stowed from From to To
	missing    map[string][]int         // file name to the chunks it never sends on its first pass
	badCRC     []string                 // files whose END carries a wrong CRC on their first pass
	receivers  []simOutage              // customers' receivers down, or answering a status
	lastSeq    map[string]uint64        // a satellite's last accepted command number at the start

	// Stopping.
	stopAt time.Time // the instant the rehearsal's context ends; zero is never
	cause  error     // the cause it ends with; nil is a plain cancel
}

// simKeys is uplink-keys.txt for the investors' week: the retired k-2028-07, k-2028-10,
// which every satellite holds, and k-2029-07, which activates on Wednesday 25 July.
const simKeys = `# uplink masters: id, activate, retire, secret
k-2028-07 2028-08-11T00:00:00Z 2028-10-24T00:00:00Z 75706c696e6b2074657374206d6173746572206b65792c203332206279746573
k-2028-10 2028-10-17T00:00:00Z 2029-08-01T00:00:00Z 75706c696e6b206e657874206d6173746572206b65792c203332206279746573
k-2029-07 2029-07-25T00:00:00Z - 75706c696e6b2032303239206d6173746572206b65792c203332206279746573
`

// simCustomers is customers.json: the five customers, their receivers and test secrets.
const simCustomers = `[
{"id":"apogee","name":"Apogee","webhook_url":"https://hooks.apogee.example/apogee","webhook_secret":"apogee webhook secret"},
{"id":"furrowlight","name":"Furrowlight","webhook_url":"https://hooks.furrowlight.example/apogee","webhook_secret":"furrowlight webhook secret"},
{"id":"brinewatch","name":"Brinewatch","webhook_url":"https://hooks.brinewatch.example/apogee","webhook_secret":"brinewatch webhook secret"},
{"id":"cirrine","name":"Cirrine Atmospherics","webhook_url":"https://hooks.cirrine.example/apogee","webhook_secret":"cirrine webhook secret"},
{"id":"kessick","name":"Kessick Polytechnic","webhook_url":"https://hooks.kessick.example/apogee","webhook_secret":"kessick webhook secret"}
]
`

// simRaceStations are the four stations a full-week case rehearses under -race, where the
// whole network would take too long.
var simRaceStations = []string{"YST", "TIR", "PEL", "VRL"}

// simLine is one Write on the console, with the instant it was written.
type simLine struct {
	At   time.Time
	Text string // as written, its newline included
}

// simConsole records the rehearsal's console.
type simConsole struct {
	mu    sync.Mutex
	now   func() time.Time
	lines []simLine
}

// Write keeps p as one console line.
func (c *simConsole) Write(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lines = append(c.lines, simLine{c.now(), string(p)})
	return len(p), nil
}

// simResult is everything one run left behind.
type simResult struct {
	c       simCase
	win     Window            // the window the case flies
	codes   []string          // the rehearsed stations, in table order
	dir     string            // the rehearsal's folder
	files   map[string]string // its files, by name with slashes, once a cached run has read them
	err     error             // what rehearse returned
	stuck   bool              // rehearse was still running an hour after the window's end
	ended   time.Time         // the instant rehearse returned
	console []simLine         // every Write on World.Out, in order
	world   *simWorld         // the stations, satellites and receivers, with their records
}

// lines returns the console's lines in order, each without its newline.
func (r *simResult) lines() []string {
	out := make([]string, len(r.console))
	for i, l := range r.console {
		out[i] = strings.TrimSuffix(l.Text, "\n")
	}
	return out
}

// file returns the text of a file in the rehearsal's folder, named with slashes, such as
// "plan/2029-07-23.txt", and whether it is there.
func (r *simResult) file(name string) (string, bool) {
	if r.files != nil {
		text, ok := r.files[name]
		return text, ok
	}
	data, err := os.ReadFile(filepath.Join(r.dir, filepath.FromSlash(name)))
	return string(data), err == nil
}

// fileLines returns the lines of a file in the rehearsal's folder, and whether it is
// there. An empty file has no lines.
func (r *simResult) fileLines(name string) ([]string, bool) {
	text, ok := r.file(name)
	if !ok || text == "" {
		return nil, ok
	}
	return strings.Split(strings.TrimSuffix(text, "\n"), "\n"), true
}

// names returns the names of the files in one folder of the rehearsal's folder, sorted;
// none when the folder is not there.
func (r *simResult) names(folder string) []string {
	var out []string
	if r.files != nil {
		for name := range r.files {
			if dir, base := filepath.Split(name); dir == folder+"/" {
				out = append(out, base)
			}
		}
	} else {
		entries, _ := os.ReadDir(filepath.Join(r.dir, filepath.FromSlash(folder)))
		for _, e := range entries {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out
}

// simCache holds the runs simOnce has made, for the life of the test binary.
var simCache struct {
	mu   sync.Mutex
	runs map[string]*simResult
}

// simOnce runs a case the first time a name is asked for, and returns the same run every
// later time, in this test and the ones after it: the whole network's week is flown once
// per test binary and checked from several angles. The run's files are read into memory
// and its folder removed at once.
func simOnce(t *testing.T, name string, c simCase) *simResult {
	t.Helper()
	simCache.mu.Lock()
	defer simCache.mu.Unlock()
	if r, ok := simCache.runs[name]; ok {
		return r
	}
	dir, err := os.MkdirTemp("", "rehearsal-")
	if err != nil {
		t.Fatalf("the simulated world could not make a folder: %v", err)
	}
	defer os.RemoveAll(dir)
	r := simRunIn(t, c, dir)
	r.files = map[string]string{}
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		rel, _ := filepath.Rel(dir, path)
		r.files[filepath.ToSlash(rel)] = string(data)
		return err
	})
	if err != nil {
		t.Fatalf("the simulated world could not read the rehearsal's folder: %v", err)
	}
	if simCache.runs == nil {
		simCache.runs = map[string]*simResult{}
	}
	simCache.runs[name] = r
	return r
}

// simWeekRun returns the investors' week as Ingrid scripts it (simInvestors), flown once
// per test binary.
func simWeekRun(t *testing.T) *simResult {
	t.Helper()
	return simOnce(t, "the investors' week", simInvestors())
}

// plan returns the plan Season 7's rule makes of the case's archive: the reference.
func (r *simResult) plan() *simPlan {
	return simPlanFor(r.win, r.codes)
}

// simEvent is one console line about a pass, such as
// "Thu 16:09:59 LOS thistle-2/PEL/02136 flown".
type simEvent struct {
	At     time.Time // the instant the line was written
	Stamp  string    // the instant the line gives, "Thu 16:09:59"
	Kind   string    // "AOS" or "LOS"
	PassID string
	Rest   string // what follows the pass ID: "dish 1", "flown", "failed line_lost"
}

// events returns the console's AOS and LOS lines, in the order they were written.
func (r *simResult) events() []simEvent {
	var out []simEvent
	for _, l := range r.console {
		f := strings.Fields(l.Text)
		if len(f) < 4 || (f[2] != "AOS" && f[2] != "LOS") {
			continue
		}
		out = append(out, simEvent{At: l.At, Stamp: f[0] + " " + f[1], Kind: f[2], PassID: f[3], Rest: strings.Join(f[4:], " ")})
	}
	return out
}

// outcomes returns how each pass ended, by pass ID: what its LOS line says after the pass
// ID, such as "flown" or "failed line_lost".
func (r *simResult) outcomes() map[string]string {
	out := map[string]string{}
	for _, e := range r.events() {
		if e.Kind == "LOS" {
			out[e.PassID] = e.Rest
		}
	}
	return out
}

// find returns the console lines that hold text, in order.
func (r *simResult) find(text string) []simLine {
	var out []simLine
	for _, l := range r.console {
		if strings.Contains(l.Text, text) {
			out = append(out, l)
		}
	}
	return out
}

// simRun runs rehearse once in the simulated world the case describes, in simulated time,
// and returns what happened. The clock starts at the window's first instant. If rehearse
// is still running an hour after the window's end, simRun fails the test, ends its context
// and hangs up every call, so that it can return.
func simRun(t *testing.T, c simCase) *simResult {
	t.Helper()
	return simRunIn(t, c, t.TempDir())
}

// simRunIn is simRun with the rehearsal's folder in dir.
func simRunIn(t *testing.T, c simCase, dir string) *simResult {
	t.Helper()
	if c.from.IsZero() {
		c.from = simWeek
	}
	if c.days == 0 {
		c.days = 7
	}
	res := &simResult{c: c, win: simWindow(c.from, c.days), codes: simCodes(c.stations), dir: dir}
	if err := simWriteFolder(res); err != nil {
		t.Fatalf("the simulated world could not write the rehearsal's folder: %v", err)
	}
	if c.passes == "" {
		res.plan() // planned before the bubble starts: the stations fly it
	}
	// Inside the bubble only one goroutine runs at a time anyway: one CPU spares the run
	// the cost of waking threads.
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	synctest.Test(t, func(t *testing.T) {
		shift := res.win.From.Sub(time.Now())
		now := func() time.Time { return time.Now().Add(shift).UTC() }
		console := simWatched(t, now)
		world := simNewWorld(res, now)
		res.world = world
		ctx, cancel := context.WithCancelCause(context.Background())
		w := World{Dir: res.dir, Now: now, Dial: world.dial, HTTP: world.client(), Out: console}

		done := make(chan error, 1)
		go func() { done <- rehearse(ctx, w) }()
		finished := make(chan struct{})
		if !c.stopAt.IsZero() {
			go func() {
				stop := time.NewTimer(c.stopAt.Sub(now()))
				defer stop.Stop()
				select {
				case <-stop.C:
					world.stopped(now())
					cancel(c.cause)
				case <-finished:
				}
			}()
		}
		guard := time.NewTimer(res.win.To.Add(time.Hour).Sub(now()))
		select {
		case res.err = <-done:
		case <-guard.C:
			res.stuck = true
			t.Errorf("rehearse still running at %s", now().Format("Mon 2 Jan 15:04"))
			cancel(nil)
			world.close()
			res.err = <-done
		}
		guard.Stop()
		res.ended = now()
		close(finished)
		cancel(nil)
		world.close()
		// Let anything the rehearsal left behind run out: a timer, a request's timeout.
		time.Sleep(max(10*time.Minute, res.win.To.Add(time.Hour).Sub(now())))
		synctest.Wait()
		console.mu.Lock()
		res.console = slices.Clone(console.lines)
		console.mu.Unlock()
	})
	return res
}

// simWriteFolder writes the rehearsal's folder as the case describes it.
func simWriteFolder(r *simResult) error {
	c := r.c
	write := func(name, text string, mode os.FileMode) error {
		return os.WriteFile(filepath.Join(r.dir, name), []byte(text), mode)
	}
	switch c.settings {
	case "-":
	case "":
		set := map[string]any{"week": c.from.Format("2006-01-02"), "days": c.days}
		if len(c.stations) > 0 {
			set["stations"] = c.stations
		}
		data, err := json.Marshal(set)
		if err != nil {
			return err
		}
		if err := write("rehearsal.json", string(data)+"\n", 0o644); err != nil {
			return err
		}
	default:
		if err := write("rehearsal.json", c.settings, 0o644); err != nil {
			return err
		}
	}
	switch c.passes {
	case "-":
	case "":
		stations := r.codes
		if len(c.archive) > 0 {
			stations = c.archive
		}
		if err := os.WriteFile(filepath.Join(r.dir, "passes.txt"), simArchiveText(r.win, stations), 0o644); err != nil {
			return err
		}
	default:
		if err := write("passes.txt", c.passes, 0o644); err != nil {
			return err
		}
	}
	switch {
	case c.bookingsText != "":
		if err := write("bookings.json", c.bookingsText, 0o644); err != nil {
			return err
		}
	case c.bookings != nil:
		if err := write("bookings.json", simBookingsText(c.bookings), 0o644); err != nil {
			return err
		}
	}
	if c.keys != "-" {
		if err := write("uplink-keys.txt", cmp.Or(c.keys, simKeys), 0o600); err != nil {
			return err
		}
	}
	if c.customers != "-" {
		if err := write("customers.json", cmp.Or(c.customers, simCustomers), 0o600); err != nil {
			return err
		}
	}
	return nil
}

// simBookingsText writes bookings as bookingd exports them: one JSON array, a booking a
// line.
func simBookingsText(bookings []Booking) string {
	var b strings.Builder
	b.WriteString("[")
	for i, bk := range bookings {
		data, err := json.Marshal(bk)
		if err != nil {
			panic(err) // a Booking always marshals
		}
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString("\n")
		b.Write(data)
	}
	b.WriteString("\n]\n")
	return b.String()
}

// simWeekBookings returns the bookings customers hold in a window at the given stations
// (none is all twenty). For a window inside the investors' week they are the week's own,
// bk-8001 to bk-9493, numbered in the whole network's plan order; under -race, and for any
// other window, they are numbered from bk-8001 in the order of the window's own plan, so
// that no test has to plan the whole network there.
func simWeekBookings(win Window, stations []string) []Booking {
	week := simWindow(simWeek, 7)
	if raceOn || win.From.Before(week.From) || win.To.After(week.To) {
		return simBook(simPlanFor(win, stations).contacts)
	}
	codes := simCodes(stations)
	var out []Booking
	for _, b := range simBook(simPlanFor(week, nil).contacts) {
		if slices.Contains(codes, b.Station) && !b.AOS.Before(win.From) && b.AOS.Before(win.To) {
			out = append(out, b)
		}
	}
	return out
}

// simSecondRun returns Yusra's second run: Thursday 26 July at Ystmark, Cape Solvane, Ilha
// Calmaria, Tirrin Rock and Quindle Island, the week's bookings there with her nine moves
// made, and the same storm at Ystmark.
func simSecondRun() simCase {
	thu := simJuly(26, 0, 0, 0)
	c := simCase{from: thu, days: 1, stations: []string{"YST", "CSL", "ICM", "TIR", "QND"}}
	moved, err := simMoved(simWeekBookings(simWindow(simWeek, 7), nil), simPlanFor(simWindow(simWeek, 7), nil))
	if err != nil {
		panic(err) // the moves name contacts of the week's plan
	}
	for _, b := range moved {
		if slices.Contains(c.stations, b.Station) && !b.AOS.Before(thu) && b.AOS.Before(thu.Add(24*time.Hour)) {
			c.bookings = append(c.bookings, b)
		}
	}
	simStorm(&c)
	return c
}

// simInvestors returns the investors' week as Ingrid scripts it: the whole network for
// seven days, its rehearsal.json naming the week alone (the four stations of
// simRaceStations under -race), the week's bookings, keys and customers, and the
// script's trouble: Varsk Island's flipped bytes, Furrow-K's three
// lost chunks on Tuesday, Cirrine's receiver restarting at 09:32 that day, and Thursday at
// Ystmark, where the fiber dies at 03:31 and the wind stows both dishes from 10:50 to 12:00.
func simInvestors() simCase {
	c := simCase{from: simWeek, days: 7, settings: `{"week":"2029-07-23"}` + "\n"}
	if raceOn {
		c.stations = simRaceStations
		c.settings = `{"week":"2029-07-23","stations":["YST","TIR","PEL","VRL"]}` + "\n"
	}
	c.bookings = simWeekBookings(simWindow(c.from, c.days), c.stations)
	simStorm(&c)
	c.flipEvery = map[string]int{"VSK": 400}
	c.missing = map[string][]int{"furrow-k-20290724-0740.tif": {21, 22, 23}}
	c.receivers = []simOutage{{Customer: "cirrine", From: simJuly(24, 9, 32, 0), To: simJuly(24, 9, 33, 0)}}
	return c
}

// simStorm adds Thursday 26 July's storm at Ystmark to a case: the fiber dies at 03:31 and
// the backhaul is up at 03:41, slow until the fiber is back at 19:40 (0.35 s more each way,
// pieces of at most 64 bytes, 0.25 s apart); the wind stows both dishes from 10:50 to 12:00.
func simStorm(c *simCase) {
	c.quiet = append(c.quiet, simSpan{Station: "YST", From: simJuly(26, 3, 31, 0), To: simJuly(26, 3, 41, 0)})
	c.backhaul = append(c.backhaul, simSlow{Station: "YST", From: simJuly(26, 3, 41, 0), To: simJuly(26, 19, 40, 0),
		Delay: 350 * time.Millisecond, Piece: 64, Gap: 250 * time.Millisecond})
	c.stows = append(c.stows, Stow{Station: "YST", From: simJuly(26, 10, 50, 0), To: simJuly(26, 12, 0, 0), Reason: "wind lock"})
}

// simMoves are Yusra's moves of Friday 20 July: each of the nine bookings the storm cost,
// with the pass ID of its Furrow's next unbooked contact that Thursday at another station
// of the second run (Furrow-V's next, at Saltreach Island, would need that station in the
// run, and a Thursday-only plan there starts its train of Furrows differently from the
// week's and loses bk-8669).
var simMoves = map[string]string{
	"bk-8743": "furrow-h/CSL/00437",
	"bk-8746": "furrow-g/ICM/00436",
	"bk-8750": "furrow-e/CSL/00436",
	"bk-8753": "furrow-z/CSL/00437",
	"bk-8755": "furrow-y/ICM/00437",
	"bk-8757": "furrow-w/CSL/00437",
	"bk-8759": "furrow-v/QND/00438",
	"bk-8762": "furrow-t/TIR/00437",
	"bk-8763": "furrow-s/QND/00438",
}

// simMoved returns bookings with Yusra's moves made: each moved booking keeps its ID and
// takes the station, AOS and LOS of its new contact, which plan must hold.
func simMoved(bookings []Booking, plan *simPlan) ([]Booking, error) {
	out := slices.Clone(bookings)
	for i, b := range out {
		id, ok := simMoves[b.ID]
		if !ok {
			continue
		}
		c, ok := plan.byID[id]
		if !ok {
			return nil, fmt.Errorf("the plan has no contact %s to move %s to", id, b.ID)
		}
		out[i].Station, out[i].AOS, out[i].LOS = c.Station, c.AOS, c.LOS
	}
	return out, nil
}

// passOfDial returns the pass a call was made for: the one its HELLO named, or else the
// contact at the call's station in progress when it was dialed.
func (cl *simCall) passOfDial(r *simResult) string {
	if cl.passID != "" {
		return cl.passID
	}
	for _, c := range r.plan().contacts {
		if c.Station == cl.station && !cl.dialed.Before(c.AOS) && cl.dialed.Before(c.LOS) {
			return c.ID
		}
	}
	return ""
}

// simWatching is the run simWatch watches, from the start of the test binary on.
var simWatching = simStartWatch()

// simWatcher holds the run simWatch watches: its test, and its console, nil between runs.
type simWatcher struct {
	mu      sync.Mutex
	t       *testing.T
	console *simConsole
}

// simStartWatch starts simWatch outside every bubble.
func simStartWatch() *simWatcher {
	w := &simWatcher{}
	go simWatch(w)
	return w
}

// simWatched returns a run's console, and has simWatch watch the run until its bubble's
// test t ends.
func simWatched(t *testing.T, now func() time.Time) *simConsole {
	c := &simConsole{now: now}
	simWatching.mu.Lock()
	simWatching.t, simWatching.console = t, c
	simWatching.mu.Unlock()
	t.Cleanup(func() {
		simWatching.mu.Lock()
		if simWatching.console == c {
			simWatching.t, simWatching.console = nil, nil
		}
		simWatching.mu.Unlock()
	})
	return c
}

// simWatch watches each run from outside its bubble, in real time, for the one hang the
// guard can't end. The bubble's clock moves only while every goroutine in it waits on
// something the bubble can see, and a sync.Mutex isn't one: a mutex held while its holder
// waits for simulated time stops the clock for good, until go test's timeout prints every
// goroutine. When a run's console has been quiet for half a second, simWatch looks at
// every goroutine; if twice in a row none in the bubble can run and one waits for a mutex,
// it fails the run's test and ends the test binary.
func simWatch(w *simWatcher) {
	var watched *simConsole
	var began time.Time
	seen, quiet, stuck := 0, 0, 0
	for range time.Tick(250 * time.Millisecond) {
		w.mu.Lock()
		t, c := w.t, w.console
		w.mu.Unlock()
		if c == nil {
			watched = nil
			continue
		}
		c.mu.Lock()
		n, last := len(c.lines), time.Time{}
		if n > 0 {
			last = c.lines[n-1].At
		}
		c.mu.Unlock()
		if c != watched || n != seen {
			if c != watched {
				watched, began = c, time.Now()
			}
			seen, quiet, stuck = n, 0, 0
			continue
		}
		if quiet++; quiet < 2 {
			continue
		}
		what, ok := simMutexStuck()
		if !ok {
			quiet, stuck = 0, 0
			continue
		}
		if stuck++; stuck < 2 {
			continue
		}
		after := ""
		if n > 0 {
			after = " after " + last.Format("Mon 15:04:05")
		}
		simFail(t, began, fmt.Sprintf("simulated time stopped%s while a goroutine waited for %s: a mutex held across a wait stops the clock", after, what))
	}
}

// simMutexStuck reports whether the bubble can never move again: no goroutine in it can
// run, and one waits for a sync.Mutex or sync.RWMutex. It also says what the first such
// goroutine waits for and where the rehearsal asked for it: "a sync.Mutex at line.go:152".
func simMutexStuck() (string, bool) {
	buf := make([]byte, 1<<20)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	what := ""
	for g := range strings.SplitSeq(string(buf), "\n\n") {
		head, frames, _ := strings.Cut(g, "\n")
		_, state, _ := strings.Cut(head, " [") // "sync.Mutex.Lock, synctest bubble 1]:"
		switch {
		case !strings.Contains(state, "synctest bubble"), strings.Contains(state, "(durable)"):
		case strings.HasPrefix(state, "sync.Mutex.Lock"), strings.HasPrefix(state, "sync.RWMutex."):
			if what == "" {
				what = simLock(frames)
			}
		default:
			return "", false // it can run, or waits on something else
		}
	}
	return what, what != ""
}

// simLock says what lock a goroutine waits for, and where its stack leaves the runtime
// and the sync package: "a sync.Mutex at line.go:152".
func simLock(frames string) string {
	what := "a sync.Mutex"
	lines := strings.Split(frames, "\n")
	for i := 0; i+1 < len(lines); i += 2 {
		switch fn := lines[i]; {
		case strings.HasPrefix(fn, "sync.(*RWMutex)"):
			what = "a sync.RWMutex"
		case strings.HasPrefix(fn, "runtime."), strings.HasPrefix(fn, "sync."), strings.HasPrefix(fn, "internal/"):
		default:
			file, _, _ := strings.Cut(strings.TrimSpace(lines[i+1]), " +0x")
			return what + " at " + filepath.Base(file)
		}
	}
	return what
}

// simFail ends the test binary with t failed and msg its message, printed as go test
// prints a failed test: a run that can never end leaves its test no other way to end.
func simFail(t *testing.T, began time.Time, msg string) {
	var b strings.Builder
	names, took := strings.Split(t.Name(), "/"), time.Since(began).Seconds()
	for i := range names {
		fmt.Fprintf(&b, "%s--- FAIL: %s (%.2fs)\n", strings.Repeat("    ", i), strings.Join(names[:i+1], "/"), took)
	}
	_, file, line, _ := runtime.Caller(1)
	fmt.Fprintf(&b, "%s%s:%d: %s\nFAIL\n", strings.Repeat("    ", len(names)), filepath.Base(file), line, msg)
	os.Stdout.WriteString(b.String())
	os.Exit(1)
}
