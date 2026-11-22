package main

import "context"

// passKey is the context key of the pass a piece of work belongs to; no other package can make one.
type passKey struct{}

// withPass returns a copy of ctx that carries pass ID id.
func withPass(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, passKey{}, id)
}

// passFrom returns the pass ID that withPass put in ctx, and false when there is none or it is empty.
func passFrom(ctx context.Context) (string, bool) {
	id, _ := ctx.Value(passKey{}).(string)
	return id, id != ""
}
