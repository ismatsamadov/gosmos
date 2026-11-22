package main

import (
	"errors"
	"testing"
	"time"
)

// waitAckWithin calls waitAck and stops the test, saying want, if it hasn't returned within 1 s.
func waitAckWithin(t *testing.T, acks <-chan string, limit time.Duration, want string) (string, error) {
	t.Helper()
	type result struct {
		ack string
		err error
	}
	done := make(chan result, 1)
	go func() {
		ack, err := waitAck(acks, limit)
		done <- result{ack, err}
	}()
	select {
	case r := <-done:
		return r.ack, r.err
	case <-time.After(time.Second):
		t.Fatalf("waitAck(acks, %v) is still waiting after 1s; want %s", limit, want)
		return "", nil
	}
}

func TestWaitAckReturnsTheAck(t *testing.T) {
	acks := make(chan string, 1)
	acks <- "ACK pipit-2 118"
	got, err := waitAckWithin(t, acks, 9*time.Minute, "the waiting ack at once")
	if got != "ACK pipit-2 118" || err != nil {
		t.Errorf(`waitAck(acks, 9m0s) = %q, %v; want "ACK pipit-2 118", nil`, got, err)
	}
}

func TestWaitAckGivesUp(t *testing.T) {
	acks := make(chan string) // Pipit-2 never answers
	got, err := waitAckWithin(t, acks, 20*time.Millisecond, "ErrNoAck after 20ms")
	if got != "" || !errors.Is(err, ErrNoAck) {
		t.Errorf(`waitAck(silent acks, 20ms) = %q, %v; want "", ErrNoAck`, got, err)
	}
}
