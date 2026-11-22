package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

// plus8 is a zone eight hours ahead of UTC: after 16:00 UTC, its clocks already show the next day.
var plus8 = time.FixedZone("UTC+8", 8*60*60)

// fakeClock is a clock the test sets by hand, which writes its times in UTC+8.
type fakeClock struct{ now time.Time }

// Now returns the clock's time, written in UTC+8.
func (c *fakeClock) Now() time.Time { return c.now.In(plus8) }

// Set moves the clock to t, forward or back.
func (c *fakeClock) Set(t time.Time) { c.now = t }

// utc returns 4 or 5 December 2028 at hh:mm:ss and ms milliseconds, UTC.
func utc(day, hh, mm, ss, ms int) time.Time {
	return time.Date(2028, 12, day, hh, mm, ss, ms*1e6, time.UTC)
}

// newRotator returns a Rotator for uplinkd in dir, with files of at most 1 KiB, on clock.
func newRotator(t *testing.T, dir string, clock *fakeClock) *Rotator {
	r := &Rotator{Dir: dir, Program: "uplinkd", MaxBytes: 1024, Now: clock.Now}
	t.Cleanup(func() { r.Close() })
	return r
}

// record returns record i, exactly size bytes long with its newline.
func record(i, size int) []byte {
	head := fmt.Sprintf("record %03d ", i)
	return []byte(head + strings.Repeat(".", size-len(head)-1) + "\n")
}

// write writes p to r at the clock's time t, as uplinkd's handler does, in one Write.
func write(t *testing.T, r *Rotator, clock *fakeClock, at time.Time, p []byte) {
	t.Helper()
	clock.Set(at)
	if n, err := r.Write(p); n != len(p) || err != nil {
		t.Fatalf("at %s, r.Write(%q…) = %d, %v, want %d, nil", at.Format("2006-01-02 15:04:05.000"), p[:10], n, err, len(p))
	}
}

// folder returns the folder's files by name and what each holds.
func folder(t *testing.T, dir string) map[string][]byte {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	files := map[string][]byte{}
	for _, e := range entries {
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		files[e.Name()] = b
	}
	return files
}

// check compares the folder with want, file by file, naming when in every message.
func check(t *testing.T, dir, when string, want map[string][]byte) {
	t.Helper()
	got := folder(t, dir)
	names := func(m map[string][]byte) []string {
		var ns []string
		for n := range m {
			ns = append(ns, n)
		}
		slices.Sort(ns)
		return ns
	}
	if !slices.Equal(names(got), names(want)) {
		t.Fatalf("%s, the folder holds %q, want %q", when, names(got), names(want))
	}
	for _, name := range names(want) {
		if !bytes.Equal(got[name], want[name]) {
			t.Errorf("%s, %s holds %d bytes, want %d: the records of its own day, in order", when, name, len(got[name]), len(want[name]))
		}
	}
}

// join returns the records one after another.
func join(recs ...[]byte) []byte { return bytes.Join(recs, nil) }

func TestRotateAtMidnight(t *testing.T) {
	dir, clock := t.TempDir(), &fakeClock{}
	r := newRotator(t, dir, clock)
	a, b, c := record(0, 90), record(1, 90), record(2, 90)
	write(t, r, clock, utc(4, 23, 59, 59, 0), a)
	write(t, r, clock, utc(4, 23, 59, 59, 999), b)
	check(t, dir, "at 23:59:59.999 on the 4th", map[string][]byte{"uplinkd.log": join(a, b)})
	write(t, r, clock, utc(5, 0, 0, 0, 0), c)
	check(t, dir, "at 00:00:00 on the 5th", map[string][]byte{
		"uplinkd-2028-12-04.1.log": join(a, b),
		"uplinkd.log":              c,
	})
	big := record(3, 1000)
	write(t, r, clock, utc(5, 0, 0, 1, 0), big)
	check(t, dir, "after a record too big for the 5th's file", map[string][]byte{
		"uplinkd-2028-12-04.1.log": join(a, b),
		"uplinkd-2028-12-05.1.log": c,
		"uplinkd.log":              big,
	})
}

func TestEmptyFileTakesTheDay(t *testing.T) {
	dir, clock := t.TempDir(), &fakeClock{}
	old1, old2 := record(10, 1000), record(11, 1000)
	for name, data := range map[string][]byte{"uplinkd-2028-12-04.1.log": old1, "uplinkd-2028-12-04.2.log": old2, "uplinkd.log": nil} {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0640); err != nil {
			t.Fatal(err)
		}
	}
	r := newRotator(t, dir, clock)
	clock.Set(utc(4, 23, 58, 0, 0)) // uplinkd starts, its file empty after the 4th's second rotation
	if err := r.open(); err != nil {
		t.Fatal(err)
	}
	a, big := record(0, 90), record(1, 1000)
	write(t, r, clock, utc(5, 0, 0, 5, 0), a)
	check(t, dir, "at 00:00:05 on the 5th, the file empty since the 4th", map[string][]byte{
		"uplinkd-2028-12-04.1.log": old1,
		"uplinkd-2028-12-04.2.log": old2,
		"uplinkd.log":              a,
	})
	write(t, r, clock, utc(5, 0, 0, 6, 0), big)
	check(t, dir, "after the 5th's first rotation", map[string][]byte{
		"uplinkd-2028-12-04.1.log": old1,
		"uplinkd-2028-12-04.2.log": old2,
		"uplinkd-2028-12-05.1.log": a,
		"uplinkd.log":              big,
	})
}

func TestClockStepsBack(t *testing.T) {
	dir, clock := t.TempDir(), &fakeClock{}
	r := newRotator(t, dir, clock)
	a, b, c, d := record(0, 90), record(1, 90), record(2, 90), record(3, 90)
	write(t, r, clock, utc(4, 23, 59, 58, 0), a)
	write(t, r, clock, utc(5, 0, 0, 0, 500), b)
	write(t, r, clock, utc(4, 23, 59, 59, 0), c) // the clock steps back a second
	write(t, r, clock, utc(5, 0, 0, 1, 0), d)
	check(t, dir, "after the clock stepped back over midnight and on again", map[string][]byte{
		"uplinkd-2028-12-04.1.log": a,
		"uplinkd.log":              join(b, c, d),
	})
}

func TestRestartNextDay(t *testing.T) {
	dir, clock := t.TempDir(), &fakeClock{}
	old1, old2, last := record(10, 1000), record(11, 1000), record(12, 300)
	for name, data := range map[string][]byte{"uplinkd-2028-12-04.1.log": old1, "uplinkd-2028-12-04.2.log": old2, "uplinkd.log": last} {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0640); err != nil {
			t.Fatal(err)
		}
	}
	at := utc(4, 23, 50, 0, 0) // uplinkd.log last changed at 23:50 on the 4th
	if err := os.Chtimes(filepath.Join(dir, "uplinkd.log"), at, at); err != nil {
		t.Fatal(err)
	}
	r := newRotator(t, dir, clock)
	a := record(0, 90)
	write(t, r, clock, utc(5, 0, 10, 0, 0), a) // the first record after a restart at 00:10
	check(t, dir, "after a restart at 00:10 on the 5th", map[string][]byte{
		"uplinkd-2028-12-04.1.log": old1,
		"uplinkd-2028-12-04.2.log": old2,
		"uplinkd-2028-12-04.3.log": last,
		"uplinkd.log":              a,
	})
}
