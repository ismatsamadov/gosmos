package main

import "sync/atomic"

// Publisher holds the schedule published last, which the Board and the dish loaders read.
type Publisher struct {
	cur atomic.Pointer[Schedule]
}

// Current returns the schedule published last, or nil before the first.
func (p *Publisher) Current() *Schedule { return p.cur.Load() }

// Update publishes contacts as a new schedule, its Version one more than the current one's, and returns it.
func (p *Publisher) Update(contacts []Contact) *Schedule {
	s := p.cur.Load()
	if s == nil {
		s = &Schedule{}
		p.cur.Store(s)
	}
	s.Version++ // TODO: build a new Schedule instead of changing the one readers hold
	s.Contacts = append(s.Contacts[:0], contacts...)
	return s
}
