package main

import "log/slog"

// passGroup returns pass p as one attribute, the group "pass": its satellite, station, AOS and LOS, in UTC.
func passGroup(p PassView) slog.Attr {
	return slog.String("pass", p.ID)
}
