package main

// Heard is one telemetry line as a station heard it.
type Heard struct {
	Station string // "TIR"
	Line    string // "cirrine-1 9012 2029-03-07T20:45:05Z batt_v=8.02 temp_c=11.5 mode=NOMINAL"
}
