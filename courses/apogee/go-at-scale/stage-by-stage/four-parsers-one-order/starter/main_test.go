package main

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync/atomic"
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

// frames returns n of sat's telemetry lines, one every 30 seconds from at, numbered from seq.
func frames(sat string, seq int, at time.Time, n int) []string {
	lines := make([]string, n)
	for i := range lines {
		t := at.Add(time.Duration(i) * 30 * time.Second)
		lines[i] = fmt.Sprintf("%s %d %s batt_v=8.01 temp_c=10.5 mode=NOMINAL", sat, seq+i, t.Format(time.RFC3339))
	}
	return lines
}

// drill returns n lines of the drill's merged feeds: five satellites in turn, every 50th line garbled.
func drill(n int) []string {
	at := time.Date(2029, 3, 7, 8, 0, 5, 0, time.UTC)
	each := n/5 + 1
	sats := [][]string{
		frames("pipit-2", 12100, at, each), frames("pipit-3", 7800, at, each), frames("furrow-a", 31400, at, each),
		frames("furrow-b", 31380, at, each), frames("brine-3", 24100, at, each),
	}
	lines := make([]string, n)
	for i := range lines {
		lines[i] = sats[i%5][i/5]
		if i%50 == 49 {
			lines[i] = strings.Replace(lines[i], "batt_v=8.01", "batt_v=8.0l", 1)
		}
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

// parsed returns the frames parseLine makes of lines, one line after another, skipping the lines it refuses.
func parsed(lines []string) []TelemetryFrame {
	var want []TelemetryFrame
	for _, l := range lines {
		if f, err := parseLine(l); err == nil {
			want = append(want, f)
		}
	}
	return want
}

// checkOrder fails the test unless got holds want's frames, in want's order.
func checkOrder(t *testing.T, call string, got, want []TelemetryFrame) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s sent %d frames, want %d", call, len(got), len(want))
	}
	for i := range got {
		if got[i].SatID != want[i].SatID || got[i].Seq != want[i].Seq {
			t.Fatalf("%s frame %d = %s %d, want %s %d, in the order the lines arrived", call, i, got[i].SatID, got[i].Seq, want[i].SatID, want[i].Seq)
		}
	}
}

func TestParseInOrder(t *testing.T) {
	g := newGuard()
	t.Run("four lines held, then let go last-first", func(t *testing.T) {
		lines := drill(12)
		held := map[string]int{lines[0]: 0, lines[1]: 1, lines[2]: 2, lines[3]: 3, lines[4]: 4}
		var inFlight, most atomic.Int64
		arrived := make(chan int, 5)
		var release, returned [5]chan struct{}
		for i := range 5 {
			release[i], returned[i] = make(chan struct{}), make(chan struct{})
		}
		// parse counts the lines in flight, and holds each of the first five until the test lets it go.
		parse := func(line string) (TelemetryFrame, error) {
			n := inFlight.Add(1)
			for m := most.Load(); n > m && !most.CompareAndSwap(m, n); m = most.Load() {
			}
			i, isHeld := held[line]
			if isHeld {
				arrived <- i
				<-release[i]
			}
			f, err := parseLine(line)
			inFlight.Add(-1)
			if isHeld {
				close(returned[i])
			}
			return f, err
		}
		var let [5]bool // the held lines the test has let go
		letGo := func(i int) {
			if !let[i] {
				let[i] = true
				close(release[i])
			}
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer func() {
			for i := range 5 {
				letGo(i)
			}
		}()
		out := parseInOrder(ctx, feed(lines, true), 4, parse)
		for range 4 {
			select {
			case <-arrived:
			case <-g.wait():
				t.Fatalf("most lines in flight = %d, want 4: parseInOrder never parsed four lines at once", most.Load())
			}
		}
		// With four lines held, a right parseInOrder starts no fifth: 50 ms gives a wrong one time to.
		select {
		case i := <-arrived:
			t.Errorf("most lines in flight = 5, want 4: line %d began while four lines were held", i+1)
		case <-time.After(50 * time.Millisecond):
		}
		for _, i := range []int{3, 2, 1, 0, 4} { // the fourth line first, the fifth last
			letGo(i)
			select {
			case <-returned[i]:
			case <-g.wait():
				t.Fatalf("the parse of line %d never returned after the test let it go", i+1)
			}
		}
		got, closed := drain(g, out)
		if !closed {
			t.Fatalf("parseInOrder(12 lines, 4 workers) left its channel open after %d frames, want it closed after the last", len(got))
		}
		checkOrder(t, "parseInOrder(12 lines, 4 workers)", got, parsed(lines))
		if m := most.Load(); m != 4 {
			t.Errorf("most lines in flight = %d, want 4", m)
		}
	})
	t.Run("2,000 lines of the drill", func(t *testing.T) {
		lines := drill(2000)
		got, closed := drain(g, parseInOrder(context.Background(), feed(lines, true), 4, parseLine))
		if !closed {
			t.Fatalf("parseInOrder(2,000 lines, 4 workers) left its channel open after %d frames, want it closed after the last", len(got))
		}
		checkOrder(t, "parseInOrder(2,000 lines, 4 workers)", got, parsed(lines))
	})
	t.Run("a garbled line keeps its place", func(t *testing.T) {
		lines := drill(5)
		lines[2] = strings.Replace(lines[2], "batt_v=8.01", "batt_v=8.0l", 1)
		got, closed := drain(g, parseInOrder(context.Background(), feed(lines, true), 4, parseLine))
		if !closed {
			t.Fatalf("parseInOrder(5 lines, the third garbled) left its channel open after %d frames, want it closed after the last", len(got))
		}
		checkOrder(t, "parseInOrder(5 lines, the third garbled)", got, parsed(lines))
	})
	t.Run("the Board takes 10 frames and leaves", func(t *testing.T) {
		before := runtime.NumGoroutine()
		ctx, cancel := context.WithCancel(context.Background())
		out := parseInOrder(ctx, feed(drill(200), false), 4, parseLine)
		for i := range 10 {
			select {
			case <-out:
			case <-g.wait():
				cancel()
				t.Fatalf("parseInOrder(200 lines) sent %d frames, want 10 before the Board leaves", i)
			}
		}
		cancel()
		if n := goroutinesLeft(g, before); n != 0 {
			t.Errorf("goroutines left after cancel = %d, want 0", n)
		}
	})
}
