package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

// The drill's feeds of 7 March from three stations: Ystmark and Tirrin Rock both heard Cirrine-1.
var (
	ystmark = []string{
		"cirrine-1 9010 2029-03-07T20:44:05Z batt_v=8.02 temp_c=11.4 mode=NOMINAL",
		"cirrine-1 9011 2029-03-07T20:44:35Z batt_v=8.02 temp_c=11.5 mode=NOMINAL",
		"cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5 mode=NOMINAL",
		"cirrine-1 9013 2029-03-07T20:45:35Z batt_v=8.02 temp_c=11.5 mode=NOMINAL",
	}
	tirrin = []string{
		"cirrine-1 9008 2029-03-07T20:43:05Z batt_v=8.03 temp_c=11.3 mode=NOMINAL",
		"cirrine-1 9009 2029-03-07T20:43:35Z batt_v=8.02 temp_c=11.4 mode=NOMINAL",
		"cirrine-1 9010 2029-03-07T20:44:05Z batt_v=8.02 temp_c=11.4 mode=NOMINAL",
		"cirrine-1 9011 2029-03-07T20:44:35Z batt_v=8.02 temp_c=11.5 mode=NOMINAL",
		"cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5 mode=NOMINAL",
	}
	portEllery = []string{
		"pipit-1 6203 2029-03-07T05:25:12Z batt_v=7.94 temp_c=3.5 mode=NOMINAL",
		"pipit-1 6204 2029-03-07T05:25:42Z batt_v=7.93 temp_c=3.6 mode=SAFE",
		"cirrine-1 9003 2029-03-07T20:40:35Z batt_v=8.04 temp_c=11.0 mode=NOMINAL",
		"cirrine-1 9004 2029-03-07T20:41:05Z batt_v=8.O3 temp_c=11.1 mode=NOMINAL",
		"cirrine-1 9005 2029-03-07T20:41:35Z batt_v=8.03 temp_c=11.2 mode=NOMINAL",
	}
)

// frames returns n of sat's telemetry lines, one every 30 seconds from at, numbered from seq.
func frames(sat string, seq int, at time.Time, n int) []string {
	lines := make([]string, n)
	for i := range lines {
		t := at.Add(time.Duration(i) * 30 * time.Second)
		lines[i] = fmt.Sprintf("%s %d %s batt_v=8.01 temp_c=10.5 mode=NOMINAL", sat, seq+i, t.Format(time.RFC3339))
	}
	return lines
}

// text is a feed's text: its lines, each ending in a newline.
func text(lines []string) string { return strings.Join(lines, "\n") + "\n" }

// guard is how long a test waits for what a right stage does at once: one second in all.
type guard struct{ end time.Time }

// newGuard starts a guard's second.
func newGuard() guard { return guard{end: time.Now().Add(time.Second)} }

// goroutinesLeft waits until no more than before goroutines run, or the guard ends, and returns how many more run.
func goroutinesLeft(g guard, before int) int {
	end := time.Now().Add(max(time.Until(g.end), 50*time.Millisecond))
	for {
		n := runtime.NumGoroutine() - before
		if n <= 0 || time.Now().After(end) {
			return max(n, 0)
		}
		time.Sleep(time.Millisecond)
	}
}

func TestIngest(t *testing.T) {
	g := newGuard()
	t.Run("three good feeds", func(t *testing.T) {
		h, err := ingest(context.Background(), map[string]io.Reader{
			"YST": strings.NewReader(text(ystmark)),
			"TIR": strings.NewReader(text(tirrin)),
			"PEL": strings.NewReader(text(portEllery)),
		})
		if err != nil {
			t.Fatalf("ingest(YST, TIR, PEL) error = %v, want nil", err)
		}
		if len(h) != 2 || h["cirrine-1"].Seq != 9013 || h["pipit-1"].Seq != 6204 {
			t.Errorf("ingest(YST, TIR, PEL) = %d satellites, Cirrine-1 at %d, Pipit-1 at %d, want 2, 9013 and 6204", len(h), h["cirrine-1"].Seq, h["pipit-1"].Seq)
		}
	})
	t.Run("Ystmark's line drops after two lines", func(t *testing.T) {
		reset := errors.New("connection reset")
		morning := time.Date(2029, 3, 7, 5, 45, 5, 0, time.UTC)
		before := runtime.NumGoroutine()
		h, err := ingest(context.Background(), map[string]io.Reader{
			"YST": io.MultiReader(strings.NewReader(text(ystmark[:2])), iotest.ErrReader(reset)),
			"TIR": strings.NewReader(text(frames("furrow-a", 31200, morning, 300))),
			"PEL": strings.NewReader(text(frames("furrow-b", 31180, morning, 300))),
		})
		if err == nil || err.Error() != "feed YST: connection reset" {
			t.Fatalf("ingest(YST broken) error = %v, want %q", err, "feed YST: connection reset")
		}
		if !errors.Is(err, reset) {
			t.Errorf("errors.Is(ingest(YST broken) error, the reader's error) = false, want true")
		}
		if h != nil {
			t.Errorf("ingest(YST broken) = %d satellites, want nil with the error", len(h))
		}
		if n := goroutinesLeft(g, before); n != 0 {
			t.Errorf("goroutines left after cancel = %d, want 0", n)
		}
	})
	t.Run("a parent context already canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := ingest(ctx, map[string]io.Reader{"PEL": strings.NewReader(text(portEllery))}); !errors.Is(err, context.Canceled) {
			t.Errorf("ingest(canceled ctx) error = %v, want context.Canceled", err)
		}
	})
}
