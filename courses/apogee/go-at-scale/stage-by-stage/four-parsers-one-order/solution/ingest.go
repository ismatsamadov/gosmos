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
	jobs := make(chan job)
	for range workers {
		go func() {
			for j := range jobs {
				f, err := parse(j.line)
				j.res <- result{f: f, err: err} // room for one: this never waits
			}
		}()
	}
	return inOrder(ctx, dispatch(ctx, in, jobs, workers))
}

// dispatch gives each line of in a result channel, queues the channels in the lines' order and hands each line to the workers.
func dispatch(ctx context.Context, in <-chan Heard, jobs chan<- job, room int) <-chan chan result {
	queue := make(chan chan result, room)
	go func() {
		defer close(queue)
		defer close(jobs)
		for {
			h, ok := receive(ctx, in)
			if !ok {
				return
			}
			res := make(chan result, 1)
			if !send(ctx, queue, res) || !send(ctx, jobs, job{line: h.Line, res: res}) {
				return
			}
		}
	}()
	return queue
}

// inOrder sends the frame from each result channel of queue, in turn, skipping lines parse refused.
func inOrder(ctx context.Context, queue <-chan chan result) <-chan TelemetryFrame {
	out := make(chan TelemetryFrame, 16)
	go func() {
		defer close(out)
		for {
			res, ok := receive(ctx, queue)
			if !ok {
				return
			}
			r, ok := receive(ctx, res)
			if !ok {
				return
			}
			if r.err == nil && !send(ctx, out, r.f) {
				return
			}
		}
	}()
	return out
}
