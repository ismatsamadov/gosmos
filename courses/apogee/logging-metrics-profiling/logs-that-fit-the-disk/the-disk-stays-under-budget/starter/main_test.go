package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

// fakeClock is a clock the test sets by hand.
type fakeClock struct{ now time.Time }

// Now returns the clock's time.
func (c *fakeClock) Now() time.Time { return c.now }

// Set moves the clock to t.
func (c *fakeClock) Set(t time.Time) { c.now = t }

// event is one record of the replay: a carrier warning, or one of Nkem's planted failures.
type event struct {
	at     time.Time
	failed bool
}

// night returns Nkem's bench night, worse than any real one, in time order: 180,000 carrier lost
// warnings, ten a second in fifty fades of six minutes, one every 14 minutes from 18:06:23 on
// 3 December, and twelve failures she plants in one minute of the third fade, more than the
// sampler's ten.
func night() []event {
	var evs []event
	start := time.Date(2028, 12, 3, 18, 6, 23, 0, time.UTC)
	for f := range 50 {
		from := start.Add(time.Duration(f) * 14 * time.Minute)
		for i := range 3600 {
			evs = append(evs, event{at: from.Add(time.Duration(i) * 100 * time.Millisecond)})
		}
	}
	for i := range 12 {
		evs = append(evs, event{at: time.Date(2028, 12, 3, 18, 37, 4*i, 50e6, time.UTC), failed: true})
	}
	slices.SortStableFunc(evs, func(a, b event) int { return a.at.Compare(b.at) })
	return evs
}

// replay plays evs through l, setting the clock to each one's time.
func replay(l *slog.Logger, clock *fakeClock, evs []event) {
	for _, ev := range evs {
		clock.Set(ev.at)
		if ev.failed {
			l.Error("pass failed", "pass_id", "cirrine-1/LAN/16402", "err", "pass cirrine-1/LAN/16402: context deadline exceeded")
		} else {
			l.Warn("carrier lost", "margin_db", -1.5)
		}
	}
}

// seed fills dir with uplinkd's rotated files from the bench's earlier runs: one of about 60 KiB
// a day, from 17 November to 2 December.
func seed(t *testing.T, dir string) {
	t.Helper()
	for d := time.Date(2028, 11, 17, 0, 0, 0, 0, time.UTC); d.Day() != 3; d = d.AddDate(0, 0, 1) {
		day := d.Format(time.DateOnly)
		line := `{"time":"` + day + `T12:00:00Z","level":"WARN","msg":"carrier lost","margin_db":-1.5}` + "\n"
		data := bytes.Repeat([]byte(line), 60<<10/len(line))
		if err := os.WriteFile(filepath.Join(dir, rotatedName("uplinkd", day, 1)), data, 0640); err != nil {
			t.Fatal(err)
		}
	}
}

// record is one line of the folder, decoded.
type record struct {
	Time  time.Time `json:"time"`
	Level string    `json:"level"`
	Msg   string    `json:"msg"`
}

// size returns how many bytes the files in dir hold together.
func size(t *testing.T, dir string) int64 {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var n int64
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			t.Fatal(err)
		}
		n += info.Size()
	}
	return n
}

func TestNewStationLog(t *testing.T) {
	dir := t.TempDir()
	seed(t, dir)
	clock := &fakeClock{now: time.Date(2028, 12, 3, 18, 0, 0, 0, time.UTC)}
	l, closer, err := newStationLog(dir, clock.Now)
	if err != nil {
		t.Fatalf("newStationLog(dir, clock.Now) returned %v, want nil", err)
	}
	defer closer.Close()
	evs := night()
	replay(l, clock, evs)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var rotatedBytes int64
	seen := map[time.Time]bool{} // the times of the carrier warnings in the folder
	failures := 0
	for _, e := range entries {
		name := e.Name()
		info, err := e.Info()
		if err != nil {
			t.Fatal(err)
		}
		if _, _, ok := parseRotated("uplinkd", name); ok {
			rotatedBytes += info.Size()
		} else if name != "uplinkd.log" {
			t.Errorf("the folder holds %s, want only uplinkd.log and files named uplinkd-<day>.<n>.log", name)
		}
		if info.Size() > fileBytes {
			t.Errorf("%s holds %d bytes, want at most %d (64 KiB)", name, info.Size(), fileBytes)
		}
		f, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		sc := bufio.NewScanner(f)
		for n := 1; sc.Scan(); n++ {
			var r record
			if err := json.Unmarshal(sc.Bytes(), &r); err != nil {
				t.Errorf("line %d of %s isn't a whole JSON record (%v)", n, name, err)
				continue
			}
			switch {
			case r.Level == "ERROR" && r.Msg == "pass failed":
				failures++
			case r.Msg == "carrier lost":
				seen[r.Time] = true
			}
		}
		f.Close()
	}
	if rotatedBytes > totalBytes {
		t.Errorf("the rotated files hold %d bytes together, want at most %d (1 MiB)", rotatedBytes, totalBytes)
	}
	if failures != 12 {
		t.Errorf("the folder holds %d of the 12 planted failures, want every ERROR line", failures)
	}
	firsts, missing := map[time.Time]int{}, 0
	var firstMissing time.Time
	for _, ev := range evs {
		m := ev.at.Truncate(time.Minute)
		if ev.failed || firsts[m] == 10 {
			continue
		}
		firsts[m]++
		if !seen[ev.at] {
			if missing == 0 {
				firstMissing = ev.at
			}
			missing++
		}
	}
	if missing > 0 {
		t.Errorf("%d of the %d warnings that open their minutes, ten a minute, aren't in the folder at their replayed times, the first from %s; want every one kept",
			missing, 10*len(firsts), firstMissing.Format("2 January 15:04:05.0"))
	}
	for _, f := range []struct {
		name  string
		there bool
	}{{"uplinkd-2028-11-19.1.log", false}, {"uplinkd-2028-12-02.1.log", true}} {
		if _, err := os.Stat(filepath.Join(dir, f.name)); (err == nil) != f.there {
			t.Errorf("after the replay, %s is there: %t, want %t: the oldest files go first", f.name, err == nil, f.there)
		}
	}
}

func TestNewStationLogClose(t *testing.T) {
	dir := t.TempDir()
	clock := &fakeClock{now: time.Date(2028, 12, 8, 14, 0, 0, 0, time.UTC)}
	l, closer, err := newStationLog(dir, clock.Now)
	if err != nil {
		t.Fatalf("newStationLog(dir, clock.Now) returned %v, want nil", err)
	}
	l.Warn("carrier lost", "margin_db", -1.5)
	if err := closer.Close(); err != nil {
		t.Fatalf("Close returned %v, want nil", err)
	}
	before := size(t, dir)
	l.Error("pass failed", "pass_id", "cirrine-1/LAN/16402")
	if after := size(t, dir); after != before {
		t.Errorf("after Close, one more record grew the folder from %d to %d bytes, want Close to close the current file", before, after)
	}
}

func TestNewStationLogMissingFolder(t *testing.T) {
	clock := &fakeClock{now: time.Date(2028, 12, 8, 14, 0, 0, 0, time.UTC)}
	if _, _, err := newStationLog(filepath.Join(t.TempDir(), "missing"), clock.Now); err == nil {
		t.Errorf("newStationLog in a folder that doesn't exist returned no error, want one")
	}
}
