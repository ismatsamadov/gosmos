package main

import (
	"slices"
	"sync/atomic"
)

// Publisher holds the schedule published last, which the Board and the dish loaders read.
type Publisher struct {
	cur atomic.Pointer[Schedule]
}

// Current returns the schedule published last, or nil before the first.
func (p *Publisher) Current() *Schedule { return p.cur.Load() }

// Update publishes contacts as a new schedule, its Version one more than the current one's, and returns it.
func (p *Publisher) Update(contacts []Contact) *Schedule {
	var version uint64
	if cur := p.cur.Load(); cur != nil {
		version = cur.Version
	}
	s := &Schedule{Version: version + 1, Contacts: slices.Clone(contacts)}
	p.cur.Store(s)
	return s
}
