package main

import (
	"context"
	"errors"
)

// metered is a Store whose Create counts every booking against its customer's plan before keeping it.
type metered struct {
	Store
	quota *Quota
}

// Create counts the booking against the plan of the customer in ctx, then keeps it; a booking the store refuses goes back to the plan.
func (m metered) Create(ctx context.Context, b Booking) (Booking, error) {
	c, ok := customerFrom(ctx)
	if !ok {
		return Booking{}, errors.New("metered: no customer in the request's context")
	}
	if err := m.quota.Take(c.ID); err != nil {
		return Booking{}, err
	}
	b, err := m.Store.Create(ctx, b)
	if err != nil {
		return Booking{}, errors.Join(err, m.quota.Give(c.ID))
	}
	return b, nil
}

// Give takes back one booking counted this month that was never made, and saves the count.
func (q *Quota) Give(customer string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.usage.Used[customer] == 0 {
		return nil
	}
	q.usage.Used[customer]--
	return q.save()
}
