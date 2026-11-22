package main

import "sync"

// Hub hands each new view of the Board to every open stream. Only the hub closes a
// stream's channel, under its mutex and only once.
type Hub struct {
	mu     sync.Mutex
	subs   map[<-chan View]chan View // each stream's channel, by the end Subscribe hands out
	closed bool                      // after Close: no new streams
}

// NewHub returns a hub with no streams.
func NewHub() *Hub {
	return &Hub{subs: map[<-chan View]chan View{}}
}

// Subscribe returns a new channel with room for 4 views, which Publish fills.
func (h *Hub) Subscribe() <-chan View {
	ch := make(chan View, 4)
	h.mu.Lock()
	defer h.mu.Unlock()
	h.subs[ch] = ch
	return ch
}

// Unsubscribe forgets ch and closes it; a channel the hub no longer holds is left alone.
func (h *Hub) Unsubscribe(ch <-chan View) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if c, ok := h.subs[ch]; ok {
		delete(h.subs, ch)
		close(c)
	}
}

// Publish sends v to every stream.
func (h *Hub) Publish(v View) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, ch := range h.subs {
		ch <- v
	}
}

// Close closes every stream's channel and refuses new streams; closing twice is harmless.
func (h *Hub) Close() {
}

// Len returns how many streams the hub holds: /health's viewers.
func (h *Hub) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.subs)
}
