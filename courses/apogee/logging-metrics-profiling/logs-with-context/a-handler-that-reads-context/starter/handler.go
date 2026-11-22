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
	return h.Handler.Handle(ctx, r) // TODO: add the context's pass and request IDs first
}
