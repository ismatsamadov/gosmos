package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// serveEvents streams the Board to one screen as server-sent events: the current view, then
// every newer one, until the screen leaves or the hub closes the stream's channel.
func (b *Board) serveEvents(w http.ResponseWriter, r *http.Request) {
	updates := b.hub.Subscribe()
	defer b.hub.Unsubscribe(updates)
	lobby := isLobby(r)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	v := b.current()
	if lobby {
		v = v.forLobby()
	}
	if err := writeEvent(w, v); err != nil {
		return
	}
	last := v.version
	for {
		select {
		case v, ok := <-updates:
			if !ok {
				fmt.Fprint(w, "event: bye\ndata: {}\n\n")
				http.NewResponseController(w).Flush()
				return
			}
			if v.version <= last {
				continue
			}
			if lobby {
				v = v.forLobby()
			}
			if err := writeEvent(w, v); err != nil {
				return
			}
			last = v.version
		case <-r.Context().Done():
			return
		}
	}
}

// writeEvent writes v as one board event and flushes it to the screen.
func writeEvent(w http.ResponseWriter, v View) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "event: board\nid: %d\ndata: %s\n\n", v.version, data); err != nil {
		return err
	}
	return http.NewResponseController(w).Flush()
}

// serveHealth answers {"status":"ok","viewers":3}: how many streams are open now.
func (b *Board) serveHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"status\":\"ok\",\"viewers\":%d}\n", b.hub.Len())
}
