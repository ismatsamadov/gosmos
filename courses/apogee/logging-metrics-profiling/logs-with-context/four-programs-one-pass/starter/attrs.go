package main

import (
	"log/slog"
	"time"
)

// msAttr returns d as whole milliseconds under key+"_ms": ("wait", d) gives wait_ms.
func msAttr(key string, d time.Duration) slog.Attr {
	return slog.Int64(key+"_ms", d.Milliseconds())
}

// passGroup returns pass p as one attribute, the group "pass": its satellite, station, AOS and LOS, in UTC.
func passGroup(p PassView) slog.Attr {
	return slog.Group("pass",
		slog.String("sat_id", p.SatID),
		slog.String("station", p.Station),
		slog.Time("aos", p.AOS.UTC()),
		slog.Time("los", p.LOS.UTC()),
	)
}
