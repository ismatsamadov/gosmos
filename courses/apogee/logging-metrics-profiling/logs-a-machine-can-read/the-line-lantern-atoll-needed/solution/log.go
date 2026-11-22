package main

import (
	"context"
	"io"
	"log/slog"
	"time"
)

// msAttr returns d as whole milliseconds under key+"_ms": ("wait", d) gives wait_ms.
func msAttr(key string, d time.Duration) slog.Attr {
	return slog.Int64(key+"_ms", d.Milliseconds())
}

// newPassLog returns uplinkd's logger: one JSON line per record on w, each stamped by now.
func newPassLog(w io.Writer, now func() time.Time) *slog.Logger {
	return slog.New(stamp(slog.NewJSONHandler(w, nil), now))
}

// passLine writes how pass p ended, wait after its preparation began: INFO "pass flown", or ERROR "pass failed" with err.
func passLine(l *slog.Logger, p PassView, wait time.Duration, err error) {
	attrs := []slog.Attr{
		slog.String("pass_id", p.ID),
		slog.Float64("max_elev_deg", p.MaxElevDeg),
		msAttr("wait", wait),
	}
	if err != nil {
		attrs = append(attrs, slog.Any("err", err))
		l.LogAttrs(context.Background(), slog.LevelError, "pass failed", attrs...)
		return
	}
	l.LogAttrs(context.Background(), slog.LevelInfo, "pass flown", attrs...)
}
