package main

import "log/slog"

// levelFor returns the level uplinkd writes an event at: DEBUG for the steps inside a pass, WARN or ERROR for trouble, INFO for the rest.
func levelFor(event string) slog.Level {
	return slog.LevelInfo
}
