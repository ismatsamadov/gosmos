package main

import (
	"context"
	"sync"
)

// merge is a fan-in: it sends every line of every feed on one channel, closes it once every feed is done, and stops when ctx is done.
func merge(ctx context.Context, feeds ...<-chan Heard) <-chan Heard {
	out := make(chan Heard, 16)
	var wg sync.WaitGroup
	for _, feed := range feeds {
		wg.Go(func() { forward(ctx, feed, out) })
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// forward sends every line of feed on out until feed is closed or ctx is done.
func forward(ctx context.Context, feed <-chan Heard, out chan<- Heard) {
	for {
		select {
		case <-ctx.Done():
			return
		case h, ok := <-feed:
			if !ok {
				return
			}
			select {
			case <-ctx.Done():
				return
			case out <- h:
			}
		}
	}
}
