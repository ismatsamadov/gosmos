package main

import (
	"fmt"
	"slices"
	"sync"
)

// Outbox keeps the deliveries still to make, and the number of the next event, in outbox.json.
// Every change is saved at once, so a restart goes on where the last process stopped.
type Outbox struct {
	path    string
	mu      sync.Mutex
	next    int        // the number of the next event ID
	pending []Delivery // in the order they were added
}

// outboxFile is what outbox.json holds: {"next_event":118,"pending":[...]}.
type outboxFile struct {
	NextEvent int        `json:"next_event"`
	Pending   []Delivery `json:"pending"`
}

// NewEventID hands out the next event ID, such as "evt-000118", once the count is saved.
func (o *Outbox) NewEventID() (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.next++
	if err := o.save(); err != nil {
		o.next--
		return "", err
	}
	return fmt.Sprintf("evt-%06d", o.next-1), nil
}

// Put adds d to the pending deliveries, or replaces the one of the same event, and saves.
func (o *Outbox) Put(d Delivery) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	before := slices.Clone(o.pending)
	i := slices.IndexFunc(o.pending, func(p Delivery) bool { return p.Event.ID == d.Event.ID })
	if i < 0 {
		o.pending = append(o.pending, d)
	} else {
		o.pending[i] = d
	}
	if err := o.save(); err != nil {
		o.pending = before
		return err
	}
	return nil
}

// Remove drops the delivery of event id, delivered or given up, and saves.
func (o *Outbox) Remove(id string) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	before := slices.Clone(o.pending)
	o.pending = slices.DeleteFunc(o.pending, func(p Delivery) bool { return p.Event.ID == id })
	if err := o.save(); err != nil {
		o.pending = before
		return err
	}
	return nil
}

// Pending returns a copy of the deliveries still to make.
func (o *Outbox) Pending() []Delivery {
	o.mu.Lock()
	defer o.mu.Unlock()
	return slices.Clone(o.pending)
}
