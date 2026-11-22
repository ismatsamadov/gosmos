package main

// newest is the sink: it drains in and returns each satellite's newest frame.
func newest(in <-chan TelemetryFrame) Health {
	h := Health{}
	for f := range in {
		h[f.SatID] = f
	}
	return h
}
