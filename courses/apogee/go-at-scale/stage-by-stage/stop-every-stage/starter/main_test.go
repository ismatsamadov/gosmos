package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"
)

// guard is how long a test waits for what a right stage does at once: one second in all.
type guard struct{ end time.Time }

// newGuard starts a guard's second.
func newGuard() guard { return guard{end: time.Now().Add(time.Second)} }

// wait is ready when the guard's second is over, and never sooner than 50 ms from now.
func (g guard) wait() <-chan time.Time {
	return time.After(max(time.Until(g.end), 50*time.Millisecond))
}

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

// pipit3 returns n of Pipit-3's lines from the drill, one every 30 seconds, numbered from 7400.
func pipit3(n int) []string {
	at := time.Date(2029, 3, 7, 6, 3, 5, 0, time.UTC)
	lines := make([]string, n)
	for i := range lines {
		t := at.Add(time.Duration(i) * 30 * time.Second)
		lines[i] = fmt.Sprintf("pipit-3 %d %s batt_v=7.88 temp_c=6.5 mode=NOMINAL", 7400+i, t.Format(time.RFC3339))
	}
	return lines
}

// feed returns a channel holding lines as Tirrin Rock heard them, closed after them when ended is true.
func feed(lines []string, ended bool) <-chan Heard {
	ch := make(chan Heard, len(lines))
	for _, l := range lines {
		ch <- Heard{Station: "TIR", Line: l}
	}
	if ended {
		close(ch)
	}
	return ch
}

// drain receives frames from ch until it closes, and reports whether it closed before the guard ended.
func drain(g guard, ch <-chan TelemetryFrame) ([]TelemetryFrame, bool) {
	var got []TelemetryFrame
	for {
		select {
		case f, ok := <-ch:
			if !ok {
				return got, true
			}
			got = append(got, f)
		case <-g.wait():
			return got, false
		}
	}
}

func TestParse(t *testing.T) {
	g := newGuard()
	t.Run("a feed of 50 lines", func(t *testing.T) {
		got, closed := drain(g, parse(context.Background(), feed(pipit3(50), true)))
		if !closed {
			t.Fatalf("parse(ctx, 50 lines) left its channel open after %d frames, want it closed after the last", len(got))
		}
		if len(got) != 50 {
			t.Fatalf("parse(ctx, 50 lines) sent %d frames, want 50", len(got))
		}
		for i, f := range got {
			if want := uint32(7400 + i); f.Seq != want {
				t.Fatalf("parse(ctx, 50 lines) frame %d has Seq = %d, want %d", i, f.Seq, want)
			}
		}
	})
	t.Run("a garbled line", func(t *testing.T) {
		lines := pipit3(3)
		lines[1] = strings.Replace(lines[1], "batt_v=7.88", "batt_v=7.8B", 1)
		got, closed := drain(g, parse(context.Background(), feed(lines, true)))
		if !closed || len(got) != 2 || got[0].Seq != 7400 || got[1].Seq != 7402 {
			t.Errorf("parse(ctx, 3 lines, the middle one garbled) = %d frames, closed %v, want frames 7400 and 7402, closed true", len(got), closed)
		}
	})
	t.Run("the Board takes 3 frames of 200 and leaves", func(t *testing.T) {
		before := runtime.NumGoroutine()
		ctx, cancel := context.WithCancel(context.Background())
		out := parse(ctx, feed(pipit3(200), false))
		for i := range 3 {
			select {
			case <-out:
			case <-g.wait():
				cancel()
				t.Fatalf("parse(ctx, 200 lines) sent %d frames, want 3 before the Board leaves", i)
			}
		}
		cancel()
		if n := goroutinesLeft(g, before); n != 0 {
			t.Errorf("goroutines left after cancel = %d, want 0", n)
		}
	})
	t.Run("canceled before the first line", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		quiet := make(chan Heard) // a feed that sends nothing and stays open
		if got, closed := drain(g, parse(ctx, quiet)); !closed {
			t.Errorf("parse(canceled ctx, a quiet feed) left its channel open after %d frames, want it closed", len(got))
		}
	})
}
