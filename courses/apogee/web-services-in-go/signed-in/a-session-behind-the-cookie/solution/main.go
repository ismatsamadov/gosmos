package main

import "crypto/rand"

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
