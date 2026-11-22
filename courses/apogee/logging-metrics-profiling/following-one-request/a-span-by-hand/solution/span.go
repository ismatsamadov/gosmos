package main

import (
	"context"
	"fmt"
)

// Start begins a span called name: the child of the span ctx carries, or the root of a new trace.
// It returns the span and a copy of ctx that carries it.
func (t *Tracer) Start(ctx context.Context, name string) (context.Context, *Span) {
	s := &Span{Name: name}
	t.mu.Lock()
	if parent, ok := spanFrom(ctx); ok {
		s.TraceID, s.ParentID = parent.TraceID, parent.SpanID
	} else {
		s.TraceID = fmt.Sprintf("%016x", t.Rand.Uint64())
	}
	s.SpanID = fmt.Sprintf("%08x", t.Rand.Uint32())
	t.mu.Unlock()
	s.PassID, _ = passFrom(ctx)
	s.Start = t.Now()
	return withSpan(ctx, s), s
}
