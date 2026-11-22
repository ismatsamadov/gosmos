package main

import (
	"fmt"
	"net/http"
)

// serveEvents streams the Board to one screen as server-sent events, until the screen
// leaves.
func (b *Board) serveEvents(w http.ResponseWriter, r *http.Request) {
	b.watch(+1)
	defer b.watch(-1)
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// writeEvent writes v as one board event and flushes it to the screen.
func writeEvent(w http.ResponseWriter, v View) error {
	return nil
}

// watch moves the count of open streams, which /health reports, by n.
func (b *Board) watch(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.viewers += n
}

// serveHealth answers {"status":"ok","viewers":3}: how many streams are open now.
func (b *Board) serveHealth(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	n := b.viewers
	b.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"status\":\"ok\",\"viewers\":%d}\n", n)
}
