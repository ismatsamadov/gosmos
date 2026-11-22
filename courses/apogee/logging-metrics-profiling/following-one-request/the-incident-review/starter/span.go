package main

import "time"

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
