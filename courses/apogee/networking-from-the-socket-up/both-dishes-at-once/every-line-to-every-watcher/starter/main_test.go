package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

const line1 = "YST-1 TRACKING az=183.25 el=41.70 wind=64"

func TestBroadcast(t *testing.T) {
	t.Run("three watchers", func(t *testing.T) {
		w := &Watchers{Timeout: time.Second}
		screens := []*screen{watching(w), watching(w), watching(w)}
		broadcast(t, w, line1)
		for i, s := range screens {
			if got := s.got(); got != line1+"\n" {
				t.Errorf("after Broadcast(%q), watcher %d got %q, want %q", line1, i+1, got, line1+"\n")
			}
		}
		if n := w.Len(); n != 3 {
			t.Errorf("after Broadcast to three watchers, Len() = %d, want 3", n)
		}
	})
	t.Run("a watcher hangs up, then ten lines", func(t *testing.T) {
		w := &Watchers{Timeout: time.Second}
		first := new(atomic.Bool) // shared: the first screen written to has hung up, whichever it is
		var screens []*screen
		for range 4 {
			s := watching(w)
			s.first = first
			screens = append(screens, s)
		}
		var want strings.Builder
		for n := 1; n <= 10; n++ {
			line := fmt.Sprintf("%s %d", line1, n)
			broadcast(t, w, line)
			want.WriteString(line + "\n")
		}
		gone := 0
		for i, s := range screens {
			if s.hasHungUp() {
				gone++
				if !s.isClosed() {
					t.Errorf("watcher %d hung up, but its connection was never closed: Remove each watcher whose write fails", i+1)
				}
			} else if got := s.got(); got != want.String() {
				t.Errorf("after ten Broadcasts, watcher %d, still listening, got %d of the ten lines: a failed write removes that watcher, and the loop goes on to the next", i+1, strings.Count(got, "\n"))
			}
		}
		if gone != 1 {
			t.Errorf("after ten Broadcasts, %d watchers had hung up, want exactly 1: the first one written to", gone)
		}
		if n := w.Len(); n != 3 {
			t.Errorf("after ten Broadcasts with one watcher gone, Len() = %d, want 3: Remove each watcher whose write fails", n)
		}
	})
	t.Run("nobody watching", func(t *testing.T) {
		w := &Watchers{Timeout: time.Second}
		broadcast(t, w, line1)
		if n := w.Len(); n != 0 {
			t.Errorf("after Broadcast to nobody, Len() = %d, want 0", n)
		}
	})
	t.Run("a watcher joins during a broadcast", func(t *testing.T) {
		w := &Watchers{Timeout: time.Second}
		s := watching(w)
		s.during = func() { w.Add(&screen{}) } // Ingrid's screen connects while the line is being written
		broadcast(t, w, line1)
		if got := s.got(); got != line1+"\n" {
			t.Errorf("after Broadcast(%q), the watcher got %q, want %q", line1, got, line1+"\n")
		}
		if n := w.Len(); n != 2 {
			t.Errorf("after a watcher joined during Broadcast, Len() = %d, want 2", n)
		}
	})
}

// broadcast calls w.Broadcast(line), failing the test if it hasn't returned within a second.
func broadcast(t *testing.T, w *Watchers, line string) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		w.Broadcast(line)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("Broadcast(%q) has not returned after 1s: a write made under the lock keeps a joining watcher waiting, so copy the set under the lock and write outside it", line)
	}
}

// watching adds a new screen to w and returns it.
func watching(w *Watchers) *screen {
	s := &screen{}
	w.Add(s)
	return s
}

// screen is one watcher's connection as dishd sees it: it keeps whatever is written to it, or,
// once the watcher has hung up, fails every write.
type screen struct {
	first  *atomic.Bool // if set, shared by screens: the first of them written to hangs up
	during func()       // if set, runs while a write is in progress

	mu      sync.Mutex
	written strings.Builder
	hungUp  bool // every write fails, as on a connection the watcher has closed
	closed  bool
}

// got returns everything written to the screen.
func (s *screen) got() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.written.String()
}

// hasHungUp reports whether the watcher has hung up.
func (s *screen) hasHungUp() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hungUp
}

// isClosed reports whether the screen's connection was closed.
func (s *screen) isClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

// Write keeps p, or fails if the watcher has hung up.
func (s *screen) Write(p []byte) (int, error) {
	s.mu.Lock()
	if s.first != nil && s.first.CompareAndSwap(false, true) {
		s.hungUp = true
	}
	hungUp := s.hungUp
	s.mu.Unlock()
	if hungUp {
		return 0, &net.OpError{Op: "write", Net: "tcp", Err: syscall.EPIPE}
	}
	if s.during != nil {
		s.during()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.written.Write(p)
}

// Read reports that the watcher sends nothing more.
func (s *screen) Read(p []byte) (int, error) { return 0, io.EOF }

// Close marks the screen's connection closed.
func (s *screen) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return nil
}

// LocalAddr returns a loopback address.
func (s *screen) LocalAddr() net.Addr { return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)} }

// RemoteAddr returns a loopback address.
func (s *screen) RemoteAddr() net.Addr { return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)} }

// SetDeadline does nothing: a screen never makes anyone wait.
func (s *screen) SetDeadline(time.Time) error { return nil }

// SetReadDeadline does nothing.
func (s *screen) SetReadDeadline(time.Time) error { return nil }

// SetWriteDeadline does nothing.
func (s *screen) SetWriteDeadline(time.Time) error { return nil }
