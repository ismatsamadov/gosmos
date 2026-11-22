package main

import "log/slog"

// levelFor returns the level uplinkd writes an event at: DEBUG for the steps inside a pass, WARN or ERROR for trouble, INFO for the rest.
func levelFor(event string) slog.Level {
	switch event {
	case "worker taken", "lines sent", "frame received", "ack received":
		return slog.LevelDebug
	case "carrier lost", "pass skipped":
		return slog.LevelWarn
	case "pass failed":
		return slog.LevelError
	}
	return slog.LevelInfo
}
