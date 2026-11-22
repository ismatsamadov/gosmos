package main

import "log/slog"

// logStart writes the first record a program logs, which names its build.
func logStart(l *slog.Logger) {
	l.Info("program started") // TODO: add the Go release, the module and its version
}
