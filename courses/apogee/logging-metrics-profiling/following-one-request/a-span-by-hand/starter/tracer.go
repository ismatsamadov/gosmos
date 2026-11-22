package main

import (
	"context"
	"math/rand/v2"
	"slices"
	"sync"
	"time"
)

// Span is one timed step of work, recorded by hand.
type Span struct {
	TraceID  string // 16 lowercase hex digits, the same for every span of one pass
	SpanID   string // 8 lowercase hex digits
	ParentID string // the enclosing span's SpanID; "" for a root
	Name     string // "deliver", "queue", "pass", "worker", "fly"
	PassID   string // from the context (passFrom), "" without one
	Start    time.Time
	End      time.Time
}

// Tracer starts spans and keeps the finished ones. Safe for many goroutines.
type Tracer struct {
	Now  func() time.Time
	Rand *rand.Rand // math/rand/v2: trace IDs aren't secrets; used only under mu
	mu   sync.Mutex
	done []Span
}

// End sets s's End from the clock and keeps a copy of s among the finished spans.
func (t *Tracer) End(s *Span) {
	s.End = t.Now()
	t.mu.Lock()
	defer t.mu.Unlock()
	t.done = append(t.done, *s)
}

// Spans returns the finished spans, in the order they ended.
func (t *Tracer) Spans() []Span {
	t.mu.Lock()
	defer t.mu.Unlock()
	return slices.Clone(t.done)
}

// spanKey is the context key a span travels under.
type spanKey struct{}

// withSpan returns a copy of ctx that carries s.
func withSpan(ctx context.Context, s *Span) context.Context {
	return context.WithValue(ctx, spanKey{}, s)
}

// spanFrom returns the span ctx carries, if it carries one.
func spanFrom(ctx context.Context) (*Span, bool) {
	s, ok := ctx.Value(spanKey{}).(*Span)
	return s, ok && s != nil
}
