package main

import (
	"context"
	"runtime"
	"strconv"
	"testing"
	"time"
)

// codes are the drill's twenty stations in table order: the five flying today, then the fifteen new ones.
var codes = []string{"YST", "PEL", "SQT", "LAN", "CWR", "TIR", "VSK", "CAL", "BRM", "OST", "ILV", "CSL", "SRH", "CRP", "NIM", "HSK", "ICM", "QND", "PVT", "VRL"}

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

// lines returns a feed holding n lines heard at station, numbered 1 to n, closed after them when ended is true.
func lines(station string, n int, ended bool) chan Heard {
	ch := make(chan Heard, n)
	for i := 1; i <= n; i++ {
		ch <- Heard{Station: station, Line: strconv.Itoa(i)}
	}
	if ended {
		close(ch)
	}
	return ch
}

// drain receives from ch until it closes, and reports whether it closed before the guard ended.
func drain(g guard, ch <-chan Heard) ([]Heard, bool) {
	var got []Heard
	for {
		select {
		case h, ok := <-ch:
			if !ok {
				return got, true
			}
			got = append(got, h)
		case <-g.wait():
			return got, false
		}
	}
}

func TestMerge(t *testing.T) {
	g := newGuard()
	t.Run("twenty feeds of 1 to 20 lines", func(t *testing.T) {
		var feeds []<-chan Heard
		for i, code := range codes {
			feeds = append(feeds, lines(code, i+1, true))
		}
		got, closed := drain(g, merge(context.Background(), feeds...))
		if !closed {
			t.Fatalf("merge(twenty feeds) sent %d of 210 lines and left its channel open, want all 210, then closed", len(got))
		}
		if len(got) != 210 {
			t.Fatalf("merge(twenty feeds) sent %d lines, want 210", len(got))
		}
		last := map[string]int{}
		for _, h := range got {
			n, _ := strconv.Atoi(h.Line)
			if n != last[h.Station]+1 {
				t.Fatalf("merge(twenty feeds) sent %s's line %d after its line %d, want each feed's lines in order", h.Station, n, last[h.Station])
			}
			last[h.Station] = n
		}
		for i, code := range codes {
			if last[code] != i+1 {
				t.Errorf("merge(twenty feeds) sent %d of %s's lines, want %d", last[code], code, i+1)
			}
		}
	})
	t.Run("no feeds", func(t *testing.T) {
		if got, closed := drain(g, merge(context.Background())); !closed || len(got) != 0 {
			t.Errorf("merge(ctx) = %d lines, closed %v, want 0 lines, closed true", len(got), closed)
		}
	})
	t.Run("nineteen feeds done, Tirrin Rock's still open", func(t *testing.T) {
		var feeds []<-chan Heard
		var tir chan Heard
		want := 0
		for i, code := range codes {
			ch := lines(code, i+1, code != "TIR")
			if code == "TIR" {
				tir = ch
			}
			feeds = append(feeds, ch)
			want += i + 1
		}
		out := merge(context.Background(), feeds...)
		for got := 0; got < want; got++ {
			select {
			case _, ok := <-out:
				if !ok {
					t.Fatalf("merge closed its channel after %d of %d lines, while Tirrin Rock's feed was still open", got, want)
				}
			case <-g.wait():
				t.Fatalf("merge sent %d of the %d lines waiting in the feeds, then nothing", got, want)
			}
		}
		select {
		case h, ok := <-out:
			if !ok {
				t.Fatalf("merge closed its channel while Tirrin Rock's feed was still open, want it open until that feed closes")
			}
			t.Fatalf("merge sent line %q from %s, which no feed holds", h.Line, h.Station)
		default:
		}
		close(tir)
		if rest, closed := drain(g, out); !closed || len(rest) != 0 {
			t.Errorf("after Tirrin Rock's feed closed, merge sent %d more lines, closed %v, want 0 lines, closed true", len(rest), closed)
		}
	})
	t.Run("canceled after 5 lines", func(t *testing.T) {
		before := runtime.NumGoroutine()
		ctx, cancel := context.WithCancel(context.Background())
		var feeds []<-chan Heard
		for _, code := range codes {
			feeds = append(feeds, lines(code, 50, false)) // live feeds, still open
		}
		out := merge(ctx, feeds...)
		for i := range 5 {
			select {
			case <-out:
			case <-g.wait():
				cancel()
				t.Fatalf("merge(twenty live feeds) sent %d lines, want 5 before the cancel", i)
			}
		}
		cancel()
		if n := goroutinesLeft(g, before); n != 0 {
			t.Errorf("goroutines left after cancel = %d, want 0", n)
		}
	})
}
