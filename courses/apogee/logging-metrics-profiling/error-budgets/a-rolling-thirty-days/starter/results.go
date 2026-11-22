package main

import "time"

// Result is how one pass ended, as a station's uplinkd reports it to the pager.
type Result struct {
	PassID  string    // "cirrine-1/LAN/16402"
	Station string    // "LAN"
	LOS     time.Time // a result counts from its LOS
	Outcome string    // "flown", "failed" or "skipped"
	Reason  string    // "" when flown; "deadline", "no_ack", "no_worker", "radio"; "stowed" when skipped
}
