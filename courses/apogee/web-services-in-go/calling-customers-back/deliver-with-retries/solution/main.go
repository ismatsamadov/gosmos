package main

import (
	"context"
	"fmt"
)

// Deliver sends d's event to its customer, trying again on the schedule, at most 8 attempts in all.
func (s *Sender) Deliver(ctx context.Context, d *Delivery) error {
	c, ok := s.Customers[d.Customer]
	if !ok {
		return fmt.Errorf("%w: %s: no customer %q", ErrUndelivered, d.Event.ID, d.Customer)
	}
	for {
		d.Attempts++
		err := post(ctx, s.Client, c.WebhookURL, c.WebhookKey, s.Now().Unix(), d.Event)
		switch {
		case err == nil:
			return nil
		case ctx.Err() != nil:
			return ctx.Err()
		case !retryable(statusOf(err), err) || d.Attempts >= 8:
			return fmt.Errorf("%w: %s, attempt %d: %w", ErrUndelivered, d.Event.ID, d.Attempts, err)
		}
		if err := s.Wait(ctx, s.schedule(d)); err != nil {
			return err
		}
	}
}
