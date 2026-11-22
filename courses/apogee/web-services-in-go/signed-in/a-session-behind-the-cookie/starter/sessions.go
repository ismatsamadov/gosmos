package main

import (
	"sync"
	"time"
)

// session is one signed-in browser: whose it is, and when it ends unless it is used again.
type session struct {
	Customer string    // "furrowlight"
	Ends     time.Time // the first instant at which the session no longer counts
}

// Sessions keeps the signed-in browsers by session ID. It is safe for many requests at once.
type Sessions struct {
	now  func() time.Time
	mu   sync.Mutex
	byID map[string]session
}

// newSessions returns an empty session store on the clock now.
func newSessions(now func() time.Time) *Sessions {
	return &Sessions{now: now, byID: map[string]session{}}
}

// End forgets session id, so its cookie signs nobody in again.
func (s *Sessions) End(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.byID, id)
}
