package main

import "context"

// feed hands contacts to a station's dish loader through q, in order, waiting while q is full: nil once q has taken them all, or ctx's error if ctx ends first.
func feed(ctx context.Context, contacts []Contact, q chan<- Contact) error {
	// TODO: send each contact into q here, waiting while q is full
	backlog := append([]Contact(nil), contacts...)
	go func() {
		for _, c := range backlog {
			q <- c
		}
	}()
	return nil
}
