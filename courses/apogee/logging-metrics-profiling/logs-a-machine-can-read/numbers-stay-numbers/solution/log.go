package main

import (
	"context"
	"log/slog"
)

// logFlown writes one INFO record for a flown pass: its ID, highest elevation and lines sent.
func logFlown(l *slog.Logger, p PassView, lines int) {
	l.LogAttrs(context.Background(), slog.LevelInfo, "pass flown",
		slog.String("pass_id", p.ID),
		slog.Float64("max_elev_deg", p.MaxElevDeg),
		slog.Int("lines", lines),
	)
}
