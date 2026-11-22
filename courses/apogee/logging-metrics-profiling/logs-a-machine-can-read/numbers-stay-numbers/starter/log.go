package main

import (
	"context"
	"log/slog"
)

// logFlown writes one INFO record for a flown pass: its ID, highest elevation and lines sent.
func logFlown(l *slog.Logger, p PassView, lines int) {
	_ = context.Background // TODO: log "pass flown" with l.LogAttrs and typed attributes
}
