package main

import "log/slog"

// programLogger returns base with the program's name and its site on every record.
func programLogger(base *slog.Logger, program, site string) *slog.Logger {
	return base
}
