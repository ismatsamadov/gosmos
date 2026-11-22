package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

// receiveWithin calls receive and stops the test if it hasn't returned within 1 s.
func receiveWithin(t *testing.T, ctx context.Context, chunks <-chan []byte) (int, error) {
	t.Helper()
	type result struct {
		n   int
		err error
	}
	done := make(chan result, 1)
	go func() {
		n, err := receive(ctx, chunks)
		done <- result{n, err}
	}()
	select {
	case r := <-done:
		return r.n, r.err
	case <-time.After(time.Second):
		t.Fatal("receive(ctx, chunks) has not returned after 1s: it must stop once ctx is done or chunks is closed")
		return 0, nil
	}
}

func TestReceiveUntilClosed(t *testing.T) {
	t.Parallel()
	chunks := make(chan []byte, 3)
	for i := range 3 {
		chunks <- []byte(fmt.Sprintf("furrow-b chunk %d", i))
	}
	close(chunks)
	n, err := receiveWithin(t, context.Background(), chunks)
	if n != 3 || err != nil {
		t.Errorf("receive(ctx, 3 chunks then closed) = %d, %v; want 3, nil", n, err)
	}
}

func TestReceiveStopsWhenCanceled(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()                    // LOS has passed
	chunks := make(chan []byte) // still open, and nothing more will come
	n, err := receiveWithin(t, ctx, chunks)
	if n != 0 || !errors.Is(err, context.Canceled) {
		t.Errorf("receive(canceled ctx, open empty channel) = %d, %v; want 0, context.Canceled", n, err)
	}
}

func TestReceiveStopsMidPass(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	chunks := make(chan []byte) // stays open: the downlink never says it's finished
	go func() {
		chunks <- []byte("furrow-b chunk 0")
		chunks <- []byte("furrow-b chunk 1")
		time.Sleep(20 * time.Millisecond) // receive is waiting for a third chunk
		cancel()                          // LOS
	}()
	n, err := receiveWithin(t, ctx, chunks)
	if n != 2 || !errors.Is(err, context.Canceled) {
		t.Errorf("receive(ctx canceled after 2 chunks, open channel) = %d, %v; want 2, context.Canceled", n, err)
	}
}
