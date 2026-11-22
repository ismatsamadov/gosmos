package main

import "context"

// rehearse runs the week rehearsal.json describes against w and returns when every pass,
// file, delivery and report is done; when ctx ends first it stops and returns
// context.Cause(ctx).
func rehearse(ctx context.Context, w World) error {
	// TODO: read rehearsal.json in w.Dir and print what the rehearsal flies
	return nil
}
