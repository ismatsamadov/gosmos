package main

import (
	"io"
	"log/slog"
)

// newLogger returns a logger that writes records to w, one line each, without the time.
func newLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{ReplaceAttr: dropTime}))
}
