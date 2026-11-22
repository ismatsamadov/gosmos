package main

import (
	"bytes"
	"errors"
	"io"
	"slices"
	"testing"
)

// recorder keeps a copy of every write, as a dish loader would.
type recorder struct{ writes [][]byte }

// Write records a copy of p.
func (r *recorder) Write(p []byte) (int, error) {
	r.writes = append(r.writes, slices.Clone(p))
	return len(p), nil
}

// lineDown is a dish loader whose line fails on write number n.
type lineDown struct{ n, writes int }

// errLineDown is lineDown's error.
var errLineDown = errors.New("dish loader line down")

// Write fails on write number n.
func (l *lineDown) Write(p []byte) (int, error) {
	l.writes++
	if l.writes == l.n {
		return 0, errLineDown
	}
	return len(p), nil
}

// takeOverExport returns take-over day's draft by station and the
// twenty station codes in table order.
func takeOverExport() (map[string][]Contact, []string) {
	byStation := map[string][]Contact{}
	for _, c := range takeOverDraft() {
		byStation[c.Station] = append(byStation[c.Station], c)
	}
	var codes []string
	for _, s := range newNetwork().Sites {
		codes = append(codes, s.Code)
	}
	return byStation, codes
}

func TestExportStations(t *testing.T) {
	byStation, codes := takeOverExport()
	t.Run("take-over day's twenty stations", func(t *testing.T) {
		var r recorder
		if err := exportStations(&r, byStation, codes); err != nil {
			t.Fatalf("exportStations(take-over day) = %v, want nil", err)
		}
		if len(r.writes) != len(codes) {
			t.Fatalf("exportStations(take-over day) made %d writes, want %d: one for each station", len(r.writes), len(codes))
		}
		for i, code := range codes {
			var want []byte
			for _, c := range byStation[code] {
				want = appendContact(want, c)
			}
			got := bytes.SplitAfter(r.writes[i], []byte("\n"))
			lines := bytes.SplitAfter(want, []byte("\n"))
			for j := range min(len(got), len(lines)) {
				if !bytes.Equal(got[j], lines[j]) {
					t.Fatalf("write %d (%s) line %d = %q, want %q", i+1, code, j+1, got[j], lines[j])
				}
			}
			if len(got) != len(lines) {
				t.Fatalf("write %d (%s) has %d lines, want %d: %s's alone", i+1, code, len(got)-1, len(lines)-1, code)
			}
		}
	})
	t.Run("Salar Quieto's line down", func(t *testing.T) {
		l := &lineDown{n: 3}
		err := exportStations(l, byStation, codes)
		if !errors.Is(err, errLineDown) {
			t.Fatalf("exportStations with SQT's line down = %v, want an error that wraps %q", err, errLineDown)
		}
		if l.writes != 3 {
			t.Errorf("exportStations with SQT's line down made %d writes, want 3: none after the failed one", l.writes)
		}
	})
}

func TestExportStationsAllocs(t *testing.T) {
	if raceOn {
		t.Skip("allocation counts differ under -race")
	}
	byStation, codes := takeOverExport()
	var err error
	allocs := testing.AllocsPerRun(10, func() {
		err = exportStations(io.Discard, byStation, codes)
	})
	if err != nil {
		return // TestExportStations says what is wrong
	}
	if allocs > 2 {
		t.Errorf("exportStations(take-over day) made %.0f allocations, want at most 2: one buffer for every station", allocs)
	}
}
