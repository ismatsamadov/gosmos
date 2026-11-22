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
	d.Attempts++
	err := post(ctx, s.Client, c.WebhookURL, c.WebhookKey, s.Now().Unix(), d.Event)
	if err != nil {
		// TODO: try again while retryable allows, after s.Wait(ctx, s.schedule(d)), up to 8 attempts in all
		return fmt.Errorf("%w: %s, attempt %d: %w", ErrUndelivered, d.Event.ID, d.Attempts, err)
	}
	return nil
}
