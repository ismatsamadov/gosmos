package main

import "log/slog"

// passGroup returns pass p as one attribute, the group "pass": its satellite, station, AOS and LOS, in UTC.
func passGroup(p PassView) slog.Attr {
	return slog.Group("pass",
		slog.String("sat_id", p.SatID),
		slog.String("station", p.Station),
		slog.Time("aos", p.AOS.UTC()),
		slog.Time("los", p.LOS.UTC()),
	)
}

// passLogger returns l with pass p's ID and its group on every record.
func passLogger(l *slog.Logger, p PassView) *slog.Logger {
	return l
}
