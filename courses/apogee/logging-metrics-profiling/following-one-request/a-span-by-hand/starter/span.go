package main

import (
	"context"
	"fmt"
)

// Start begins a span called name: the child of the span ctx carries, or the root of a new trace.
// It returns the span and a copy of ctx that carries it.
func (t *Tracer) Start(ctx context.Context, name string) (context.Context, *Span) {
	_ = fmt.Sprintf // TODO: take or draw the IDs, set Start and PassID, and return ctx carrying the span
	return ctx, &Span{Name: name}
}
