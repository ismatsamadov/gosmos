package main

import "log/slog"

// dropTime leaves the time out of every record, so a line fits the screen.
func dropTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}
