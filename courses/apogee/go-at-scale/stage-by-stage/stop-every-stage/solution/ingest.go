package main

import "context"

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
