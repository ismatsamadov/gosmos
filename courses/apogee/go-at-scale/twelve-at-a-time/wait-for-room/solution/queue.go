package main

import "context"

// feed hands contacts to a station's dish loader through q, in order, waiting while q is full: nil once q has taken them all, or ctx's error if ctx ends first.
func feed(ctx context.Context, contacts []Contact, q chan<- Contact) error {
	for _, c := range contacts {
		select {
		case q <- c:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
