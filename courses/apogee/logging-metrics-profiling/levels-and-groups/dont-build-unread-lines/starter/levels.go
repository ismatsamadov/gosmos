package main

import (
	"context"
	"log/slog"
)

// logFrame writes a DEBUG record of one received frame, its size and dump's hex of it, calling dump only when DEBUG records are written.
func logFrame(ctx context.Context, l *slog.Logger, frame []byte, dump func([]byte) string) {
	l.DebugContext(ctx, "frame received", "bytes", len(frame), "hex", dump(frame))
}
