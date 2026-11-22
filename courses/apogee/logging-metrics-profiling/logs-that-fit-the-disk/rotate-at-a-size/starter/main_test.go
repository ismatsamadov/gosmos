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

// monday is the bench clock's time: Monday 4 December, after the disk page.
func monday() time.Time { return time.Date(2028, 12, 4, 9, 30, 0, 0, time.UTC) }

// newRotator returns a Rotator for uplinkd in a fresh folder, with files of at most 1 KiB.
func newRotator(t *testing.T) *Rotator {
	r := &Rotator{Dir: t.TempDir(), Program: "uplinkd", MaxBytes: 1024, Now: monday}
	t.Cleanup(func() { r.Close() })
	return r
}

// record returns record i, exactly size bytes long with its newline.
func record(i, size int) []byte {
	head := fmt.Sprintf("record %03d ", i)
	return []byte(head + strings.Repeat(".", size-len(head)-1) + "\n")
}

// write writes each record to r as uplinkd's handler does, one Write per record.
func write(t *testing.T, r *Rotator, recs ...[]byte) {
	t.Helper()
	for _, p := range recs {
		if n, err := r.Write(p); n != len(p) || err != nil {
			t.Fatalf("r.Write(%q…) = %d, %v, want %d, nil", p[:10], n, err, len(p))
		}
	}
}

// files returns the folder's files in the order they were written: the rotated ones by number, then
// uplinkd.log; it stops the test at any other name, or at a gap in the numbers.
func files(t *testing.T, r *Rotator) (names []string, data [][]byte) {
	t.Helper()
	entries, err := os.ReadDir(r.Dir)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, e := range entries {
		if e.Name() != "uplinkd.log" {
			n++
		}
	}
	for i := 1; i <= n; i++ {
		names = append(names, rotatedName("uplinkd", "2028-12-04", i))
	}
	names = append(names, "uplinkd.log")
	for _, e := range entries {
		if !slices.Contains(names, e.Name()) {
			t.Fatalf("the folder holds %s, want only uplinkd.log and %s to %s", e.Name(), names[0], names[len(names)-2])
		}
	}
	for _, name := range names {
		b, err := os.ReadFile(filepath.Join(r.Dir, name))
		if err != nil {
			t.Fatalf("the folder has no %s (%v), want every number from 1 to %d", name, err, n)
		}
		data = append(data, b)
	}
	return names, data
}

func TestRotateAtSize(t *testing.T) {
	r := newRotator(t)
	var recs [][]byte
	for i := range 80 {
		recs = append(recs, record(i, 60+i*53%200))
	}
	write(t, r, recs...)
	names, data := files(t, r)
	if len(names) < 10 {
		t.Fatalf("80 records of 60 to 259 bytes left %q, want rotation into a dozen or so files of at most 1 KiB", names)
	}
	if got, want := bytes.Join(data, nil), bytes.Join(recs, nil); !bytes.Equal(got, want) {
		t.Fatalf("the files, read in order, hold %d bytes, want the %d written, every record once and whole", len(got), len(want))
	}
	for i, name := range names {
		if len(data[i]) > 1024 {
			t.Errorf("%s holds %d bytes, want at most 1024", name, len(data[i]))
		}
		if i+1 < len(names) {
			next := bytes.IndexByte(data[i+1], '\n') + 1
			if len(data[i])+next <= 1024 {
				t.Errorf("%s was rotated at %d bytes, though the next record, %d bytes, fit: rotate only when it doesn't", name, len(data[i]), next)
			}
		}
	}
}

func TestRotateExactlyFull(t *testing.T) {
	r := newRotator(t)
	write(t, r, record(0, 256), record(1, 256), record(2, 256), record(3, 256))
	if names, _ := files(t, r); len(names) != 1 {
		t.Errorf("four records of 256 bytes made %d files, want uplinkd.log alone: exactly 1 KiB doesn't pass the limit", len(names))
	}
	write(t, r, record(4, 100))
	names, data := files(t, r)
	if len(names) != 2 || len(data[0]) != 1024 || len(data[1]) != 100 {
		t.Errorf("one more record of 100 bytes left %q, want the full uplinkd-2028-12-04.1.log of 1024 bytes and uplinkd.log of 100", names)
	}
}

func TestRotateLongRecord(t *testing.T) {
	r := newRotator(t)
	long := record(0, 1500) // longer than the limit, such as a frame dumped as hex
	write(t, r, long)
	if names, data := files(t, r); len(names) != 1 || len(data[0]) != 1500 {
		t.Fatalf("a 1500-byte record into an empty file left %q, want it alone in uplinkd.log, with no empty file rotated", names)
	}
	write(t, r, record(1, 100))
	names, data := files(t, r)
	if len(names) != 2 || !bytes.Equal(data[0], long) {
		t.Errorf("the next record left %q, want the long record alone in uplinkd-2028-12-04.1.log", names)
	}
	write(t, r, record(2, 1000))
	if names, _ := files(t, r); len(names) != 3 {
		t.Errorf("a 1000-byte record after 100 bytes left %q, want uplinkd-2028-12-04.2.log too", names)
	}
}
