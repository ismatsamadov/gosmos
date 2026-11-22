package main

import (
	"crypto/rand"
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

// Start begins a session for customer and returns its ID, a random name nobody can guess.
func (s *Sessions) Start(customer string) string {
	id := rand.Text()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byID[id] = session{Customer: customer, Ends: s.now().Add(sessionLife)}
	return id
}

// Lookup returns the customer of session id and renews the session, or false when there is no such session or it has ended.
func (s *Sessions) Lookup(id string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ss, ok := s.byID[id]
	if !ok {
		return "", false
	}
	now := s.now()
	if !now.Before(ss.Ends) {
		delete(s.byID, id)
		return "", false
	}
	ss.Ends = now.Add(sessionLife)
	s.byID[id] = ss
	return ss.Customer, true
}

// End forgets session id, so its cookie signs nobody in again.
func (s *Sessions) End(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.byID, id)
}
