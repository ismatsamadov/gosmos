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

// parse is a stage: it sends a frame for each line of in that parseLine reads, drops the rest, and stops when ctx is done.
func parse(ctx context.Context, in <-chan Heard) <-chan TelemetryFrame {
	out := make(chan TelemetryFrame, 16)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case h, ok := <-in:
				if !ok {
					return
				}
				f, err := parseLine(h.Line)
				if err != nil {
					continue // a line nobody can read, such as one the fog garbled
				}
				select {
				case <-ctx.Done():
					return
				case out <- f:
				}
			}
		}
	}()
	return out
}

// newest is the sink: it drains in and returns each satellite's newest frame.
func newest(in <-chan TelemetryFrame) Health {
	h := Health{}
	for f := range in {
		if old, ok := h[f.SatID]; !ok || f.Seq > old.Seq {
			h[f.SatID] = f
		}
	}
	return h
}
