package main

import (
	"log/slog"
	"time"
)

// msAttr returns d as whole milliseconds under key+"_ms": ("wait", d) gives wait_ms.
func msAttr(key string, d time.Duration) slog.Attr {
	return slog.Int64(key+"_ms", d.Milliseconds())
}
