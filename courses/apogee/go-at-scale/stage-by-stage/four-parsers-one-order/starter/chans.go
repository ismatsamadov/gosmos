package main

import "context"

// send sends v on ch, or gives up when ctx is done; it reports whether v was sent.
func send[T any](ctx context.Context, ch chan<- T, v T) bool {
	select {
	case <-ctx.Done():
		return false
	case ch <- v:
		return true
	}
}

// receive receives a value from ch, or gives up when ctx is done; ok is false once ch is closed or ctx is done.
func receive[T any](ctx context.Context, ch <-chan T) (v T, ok bool) {
	select {
	case <-ctx.Done():
	case v, ok = <-ch:
	}
	return v, ok
}
