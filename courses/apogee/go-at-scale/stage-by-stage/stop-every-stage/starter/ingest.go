package main

import "context"

// parse is a stage: it sends a frame for each line of in that parseLine reads, and drops the rest.
func parse(ctx context.Context, in <-chan Heard) <-chan TelemetryFrame {
	out := make(chan TelemetryFrame, 16)
	go func() {
		defer close(out)
		for h := range in {
			f, err := parseLine(h.Line)
			if err != nil {
				continue // a line nobody can read, such as one the fog garbled
			}
			out <- f
		}
	}()
	return out
}
