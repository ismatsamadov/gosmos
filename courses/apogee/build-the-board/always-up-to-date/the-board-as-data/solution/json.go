package main

import "net/http"

// serveJSON writes the current view as JSON.
func (b *Board) serveJSON(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	data := b.json
	b.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(data)
}
