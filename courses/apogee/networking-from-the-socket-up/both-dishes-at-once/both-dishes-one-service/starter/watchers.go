package main

import (
	"net"
	"sync"
	"time"
)

// Watchers are the connections on the watch port. Safe for many goroutines.
type Watchers struct {
	Timeout time.Duration // how long one line may take to write to one watcher: 2s
	mu      sync.Mutex
	conns   map[net.Conn]bool
	closed  bool // set by CloseAll: no watcher joins after it
}

// Add puts the watcher c in the set, or closes c at once if CloseAll has run.
func (w *Watchers) Add(c net.Conn) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		c.Close()
		return
	}
	if w.conns == nil {
		w.conns = make(map[net.Conn]bool)
	}
	w.conns[c] = true
}

// Remove takes the watcher c out of the set and closes it.
func (w *Watchers) Remove(c net.Conn) {
	w.mu.Lock()
	delete(w.conns, c)
	w.mu.Unlock()
	c.Close()
}

// CloseAll closes every watcher and empties the set for good: a watcher added after it is
// closed at once.
func (w *Watchers) CloseAll() {
	w.mu.Lock()
	conns := w.conns
	w.conns = nil
	w.closed = true
	w.mu.Unlock()
	for c := range conns {
		c.Close()
	}
}

// Len returns how many watchers are in the set.
func (w *Watchers) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.conns)
}
