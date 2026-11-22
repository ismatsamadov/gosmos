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
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			out <- Heard{Station: station, Line: sc.Text()}
		}
	}()
	return out
}
