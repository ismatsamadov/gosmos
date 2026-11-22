package main

import (
	"log"
	"log/slog"
)

// init makes log write as uplinkd's old lines did, its clock stopped at 02:46:24.
func init() {
	log.SetFlags(0)
	log.SetPrefix("2028/11/08 02:46:24 ")
}

// dropTime leaves the time out of every record, so a line fits the screen.
func dropTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}
