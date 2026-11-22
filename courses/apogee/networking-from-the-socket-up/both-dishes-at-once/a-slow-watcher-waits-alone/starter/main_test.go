package main

import (
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

const (
	line1 = "YST-1 TRACKING az=183.25 el=41.70 wind=64"
	line2 = "YST-2 IDLE az=180.00 el=45.00 wind=64"
)

// The test runs in a synctest bubble, whose clock is fake: w.Timeout keeps the 2s it has at
// Port Ellery, and the screens that read take each line a byte every 30ms, 1.26s for a whole
// one, all in no real time. A screen that slow keeps up with a deadline of its own for each
// write, but two of them can't share one, whatever order the set gives them.

func TestBroadcast(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		w := &Watchers{Timeout: 2 * time.Second}
		readers := []*watcher{newWatcher(t, w, true), newWatcher(t, w, true)}
		asleep := newWatcher(t, w, false) // Tomás's laptop: connected, never reading

		if !broadcastWithin(w, line1, time.Minute) {
			asleep.dishd.Close() // let the stuck write go, so Broadcast can end
			t.Fatalf("Broadcast(%q) is still writing to a watcher that never reads after 1m0s, want it to give up after w.Timeout, 2s: set a write deadline before each write", line1)
		}
		for i, r := range readers {
			if got := r.next(); got != line1+"\n" {
				t.Errorf("after Broadcast(%q), reading watcher %d read %q, want %q: give each write its own deadline, set just before it", line1, i+1, got, line1+"\n")
			}
		}
		if n := w.Len(); n != 2 {
			t.Errorf("after Broadcast(%q) with one watcher asleep, Len() = %d, want 2: remove the watcher whose write timed out, and only that one", line1, n)
		}
		asleep.end.SetReadDeadline(time.Now().Add(time.Second))
		if _, err := asleep.end.Read(make([]byte, 64)); !errors.Is(err, io.EOF) {
			t.Errorf("the sleeping watcher's next read = %v, want io.EOF: Remove closes the watcher it takes out", err)
		}

		if !broadcastWithin(w, line2, time.Minute) {
			t.Fatalf("Broadcast(%q) after the sleeper was removed has not returned after 1m0s", line2)
		}
		for i, r := range readers {
			if got := r.next(); got != line2+"\n" {
				t.Errorf("after Broadcast(%q), reading watcher %d read %q, want %q: each write needs a deadline of its own", line2, i+1, got, line2+"\n")
			}
		}
	})
}

// watcher is one watcher on a pipe: dishd's end is in the set, and the test holds the other.
type watcher struct {
	dishd, end net.Conn
	lines      chan string // every whole line the watcher reads, if it reads
}

// newWatcher adds a watcher to w. If reading is true, its end reads every line it is sent, a
// byte every 30ms; otherwise it never reads. Both ends close when the test ends.
func newWatcher(t *testing.T, w *Watchers, reading bool) *watcher {
	t.Helper()
	dishd, end := net.Pipe()
	wt := &watcher{dishd: dishd, end: end, lines: make(chan string, 16)}
	var wg sync.WaitGroup
	t.Cleanup(func() {
		dishd.Close()
		end.Close()
		wg.Wait()
	})
	if reading {
		wg.Go(func() {
			var line []byte
			b := make([]byte, 1)
			for {
				time.Sleep(30 * time.Millisecond)
				if _, err := end.Read(b); err != nil {
					return
				}
				line = append(line, b[0])
				if b[0] == '\n' {
					wt.lines <- string(line)
					line = nil
				}
			}
		})
	}
	w.Add(dishd)
	return wt
}

// next returns the next whole line the watcher read, or "" if none comes within a minute.
func (wt *watcher) next() string {
	select {
	case line := <-wt.lines:
		return line
	case <-time.After(time.Minute):
		return ""
	}
}

// broadcastWithin calls w.Broadcast(line) and reports whether it returned within limit. If it
// didn't, it keeps waiting in the background until the test lets its write go.
func broadcastWithin(w *Watchers, line string, limit time.Duration) bool {
	done := make(chan struct{})
	go func() {
		w.Broadcast(line)
		close(done)
	}()
	select {
	case <-done:
		return true
	case <-time.After(limit):
		return false
	}
}
