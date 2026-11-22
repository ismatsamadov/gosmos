package main

import (
	"log"
	"log/slog"
)

// accessLog returns the logger bookingd's access lines go to: each line becomes one INFO record for h.
func accessLog(h slog.Handler) *log.Logger {
	return slog.NewLogLogger(h, slog.LevelInfo)
}
