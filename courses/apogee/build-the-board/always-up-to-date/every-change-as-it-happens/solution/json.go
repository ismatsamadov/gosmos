package main

import (
	"fmt"
	"net/http"
)

// serveJSON writes the current view as JSON, named by its version in an ETag, or answers
// 304 Not Modified when the screen already has that version.
func (b *Board) serveJSON(w http.ResponseWriter, r *http.Request) {
	b.mu.Lock()
	v, data := b.view, b.json
	b.mu.Unlock()
	tag := fmt.Sprintf(`"v%d"`, v.version)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("ETag", tag)
	if r.Header.Get("If-None-Match") == tag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
