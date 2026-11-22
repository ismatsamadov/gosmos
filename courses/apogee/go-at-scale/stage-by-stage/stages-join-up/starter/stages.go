package main

import (
	"bufio"
	"context"
	"io"
)

// read is the first stage: it sends each line of r, heard at station, and closes its channel at the end of r.
func read(ctx context.Context, station string, r io.Reader) <-chan Heard {
	out := make(chan Heard, 16)
	go func() {
		defer close(out)
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			out <- Heard{Station: station, Line: sc.Text()}
		}
	}()
	return out
}

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
