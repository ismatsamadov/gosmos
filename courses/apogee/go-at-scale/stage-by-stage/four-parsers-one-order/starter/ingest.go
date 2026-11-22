package main

import "context"

// result is what parse made of one line.
type result struct {
	f   TelemetryFrame
	err error
}

// job is one line to parse and the channel its result goes into, which has room for one.
type job struct {
	line string
	res  chan result
}

// parseInOrder parses the lines of in in up to workers goroutines and sends the frames in the order their lines arrived, dropping lines parse refuses.
func parseInOrder(ctx context.Context, in <-chan Heard, workers int, parse func(string) (TelemetryFrame, error)) <-chan TelemetryFrame {
	out := make(chan TelemetryFrame, 16)
	go func() { // TODO: parse in workers goroutines, and send the frames in the lines' order
		defer close(out)
		for {
			h, ok := receive(ctx, in)
			if !ok {
				return
			}
			if f, err := parse(h.Line); err == nil && !send(ctx, out, f) {
				return
			}
		}
	}()
	return out
}
