package main

import (
	"context"
	"log/slog"
	"time"
)

// dropTime leaves the time out of every record, so a line fits the screen.
func dropTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}

// stamp returns h with every record's time taken from now, in UTC; WithAttrs and WithGroup keep it.
func stamp(h slog.Handler, now func() time.Time) slog.Handler {
	return stamped{h, now}
}

// stamped is a handler that gives each record the clock's time before passing it on.
type stamped struct {
	slog.Handler
	now func() time.Time
}

// Handle sets r's time from the clock and passes r on.
func (s stamped) Handle(ctx context.Context, r slog.Record) error {
	r.Time = s.now().UTC()
	return s.Handler.Handle(ctx, r)
}

// WithAttrs returns the handler with attrs added, still stamped.
func (s stamped) WithAttrs(attrs []slog.Attr) slog.Handler {
	return stamped{s.Handler.WithAttrs(attrs), s.now}
}

// WithGroup returns the handler with the group opened, still stamped.
func (s stamped) WithGroup(name string) slog.Handler {
	return stamped{s.Handler.WithGroup(name), s.now}
}
