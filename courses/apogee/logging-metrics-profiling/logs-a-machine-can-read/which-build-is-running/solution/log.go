package main

import (
	"log/slog"
	"runtime/debug"
)

// logStart writes the first record a program logs, which names its build.
func logStart(l *slog.Logger) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		l.Info("program started")
		return
	}
	l.Info("program started",
		"go_version", bi.GoVersion,
		"module", bi.Main.Path,
		"module_version", bi.Main.Version)
}
