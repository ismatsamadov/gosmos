package main

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

// now is Monday 4 December, 20:00 UTC, as a clock in UTC+8 writes it: already the 5th there.
var now = time.Date(2028, 12, 4, 20, 0, 0, 0, time.UTC).In(time.FixedZone("UTC+8", 8*60*60))

// fill makes each named file in dir, size bytes long.
func fill(t *testing.T, dir string, size int, names ...string) {
	t.Helper()
	for _, name := range names {
		if err := os.WriteFile(filepath.Join(dir, name), bytes.Repeat([]byte("."), size), 0640); err != nil {
			t.Fatal(err)
		}
	}
}

// left returns the names of the files in dir, in the order os.ReadDir gives them.
func left(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

// others are files prune must never touch, however old or large: the current file, another
// program's file, and names that don't parse.
var others = []string{"uplinkd.log", "bookingd-2028-11-01.1.log", "uplinkd-2028-11-02.1.log.gz", "uplinkd-old.log", "notes.txt"}

func TestPruneByAge(t *testing.T) {
	dir := t.TempDir()
	fill(t, dir, 100, "uplinkd-2028-11-19.1.log", "uplinkd-2028-11-19.2.log", "uplinkd-2028-11-20.1.log", "uplinkd-2028-12-04.1.log")
	fill(t, dir, 100, others...)
	if err := prune(dir, "uplinkd", now, 14, 1<<30); err != nil {
		t.Fatalf("prune(dir, \"uplinkd\", now, 14, 1 GiB) returned %v, want nil", err)
	}
	want := []string{"bookingd-2028-11-01.1.log", "notes.txt", "uplinkd-2028-11-02.1.log.gz", "uplinkd-2028-11-20.1.log", "uplinkd-2028-12-04.1.log", "uplinkd-old.log", "uplinkd.log"}
	if got := left(t, dir); !slices.Equal(got, want) {
		t.Errorf("on 4 December, keeping 14 days, prune left %q, want %q: 20 November is 14 days back and stays, 19 November is 15 and goes", got, want)
	}
}

func TestPruneBySize(t *testing.T) {
	dir := t.TempDir()
	fill(t, dir, 300, "uplinkd-2028-11-30.1.log", "uplinkd-2028-12-01.2.log", "uplinkd-2028-12-01.9.log",
		"uplinkd-2028-12-01.10.log", "uplinkd-2028-12-03.1.log", "uplinkd-2028-12-04.1.log")
	fill(t, dir, 5000, others...)
	if err := prune(dir, "uplinkd", now, 14, 1000); err != nil {
		t.Fatalf("prune(dir, \"uplinkd\", now, 14, 1000) returned %v, want nil", err)
	}
	want := []string{"bookingd-2028-11-01.1.log", "notes.txt", "uplinkd-2028-11-02.1.log.gz", "uplinkd-2028-12-01.10.log",
		"uplinkd-2028-12-03.1.log", "uplinkd-2028-12-04.1.log", "uplinkd-old.log", "uplinkd.log"}
	if got := left(t, dir); !slices.Equal(got, want) {
		t.Errorf("six rotated files of 300 bytes in a budget of 1000: prune left %q, want %q: the oldest go first, by day, then number, .9 before .10", got, want)
	}
}

func TestPruneWithinBudget(t *testing.T) {
	dir := t.TempDir()
	fill(t, dir, 300, "uplinkd-2028-12-01.1.log", "uplinkd-2028-12-02.1.log", "uplinkd-2028-12-03.1.log")
	fill(t, dir, 5000, others...)
	before := left(t, dir)
	if err := prune(dir, "uplinkd", now, 14, 900); err != nil {
		t.Fatalf("prune(dir, \"uplinkd\", now, 14, 900) returned %v, want nil", err)
	}
	if got := left(t, dir); !slices.Equal(got, before) {
		t.Errorf("three rotated files of 300 bytes in a budget of 900: prune left %q, want every file kept: they fit exactly", got)
	}
}

func TestPruneUnreadableFolder(t *testing.T) {
	if err := prune(filepath.Join(t.TempDir(), "missing"), "uplinkd", now, 14, 1<<30); err == nil {
		t.Errorf("prune on a folder that doesn't exist returned nil, want its error")
	}
}
