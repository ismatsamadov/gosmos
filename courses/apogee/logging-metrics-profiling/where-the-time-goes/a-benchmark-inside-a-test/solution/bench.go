package main

import "testing"

// benchPlan measures planDay on the inputs from benchInputs, built once.
func benchPlan(b *testing.B) {
	in := benchInputs()
	for b.Loop() {
		planDay(in.sats, in.els, in.stations, in.bookings, in.from, in.to)
	}
}
