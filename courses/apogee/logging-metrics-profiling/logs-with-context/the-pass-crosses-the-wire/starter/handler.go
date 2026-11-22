package main

import (
	"context"
	"log/slog"
)

// passHandler adds the pass and request IDs in a record's context to the record, then hands it to the handler inside.
type passHandler struct {
	slog.Handler
}

// Handle adds pass_id and request_id from ctx, when it has them, after r's own attributes, and passes r on.
func (h passHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := passFrom(ctx); ok {
		r.AddAttrs(slog.String("pass_id", id))
	}
	if id := requestIDFrom(ctx); id != "" {
		r.AddAttrs(slog.String("request_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

// WithAttrs returns the handler with attrs added, still adding the context's IDs.
func (h passHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return passHandler{h.Handler.WithAttrs(attrs)}
}

// WithGroup returns the handler with the group opened, still adding the context's IDs.
func (h passHandler) WithGroup(name string) slog.Handler {
	return passHandler{h.Handler.WithGroup(name)}
}
