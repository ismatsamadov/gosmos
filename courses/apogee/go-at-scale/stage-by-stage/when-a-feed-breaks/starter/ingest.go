package main

import (
	"context"
	"io"
)

// ingest runs tlmd's pipeline over every feed and returns each satellite's newest frame, or the error that broke a feed.
func ingest(ctx context.Context, feeds map[string]io.Reader) (Health, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var lines []<-chan Heard
	for station, r := range feeds {
		lines = append(lines, readFeed(ctx, func(error) { cancel() }, station, r))
	}
	return newest(parse(ctx, merge(ctx, lines...))), nil
}
