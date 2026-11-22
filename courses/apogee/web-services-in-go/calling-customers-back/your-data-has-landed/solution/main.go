package main

import (
	"context"
	"fmt"
)

// Landed tells the booking's customer, with a signed data.landed event, that its file has landed.
func (n *Notifier) Landed(ctx context.Context, l Landing) error {
	b, err := n.Store.Get(ctx, l.BookingID)
	if err != nil {
		return err
	}
	c, ok := n.owner(b.SatID)
	if !ok {
		return fmt.Errorf("booking %s: no customer operates %s", b.ID, b.SatID)
	}
	ev := Event{
		ID:        fmt.Sprintf("evt-%06d", n.events.Add(1)),
		Type:      "data.landed",
		BookingID: b.ID,
		SatID:     b.SatID,
		Station:   b.Station,
		File:      l.File,
		Bytes:     l.Bytes,
		LandedAt:  l.At,
	}
	return n.Sender.Deliver(ctx, &Delivery{Customer: c.ID, Event: ev})
}
