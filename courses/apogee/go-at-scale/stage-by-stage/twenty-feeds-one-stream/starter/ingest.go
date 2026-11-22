package main

import (
	"context"
	"sync"
)

// merge is a fan-in: it sends every line of every feed on one channel, closes it once every feed is done, and stops when ctx is done.
func merge(ctx context.Context, feeds ...<-chan Heard) <-chan Heard {
	out := make(chan Heard, 16)
	var wg sync.WaitGroup // TODO: forward each feed in a goroutine of wg, then close out after wg.Wait()
	_ = &wg
	return out
}
