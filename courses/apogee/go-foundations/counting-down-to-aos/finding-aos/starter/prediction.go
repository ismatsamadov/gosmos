package main

// elevationAt returns Pipit-3's predicted elevation in degrees at a minute of the prediction window.
func elevationAt(minute int) float64 {
	d := float64(minute - 7)
	return 64 - 1.6*d*d
}
