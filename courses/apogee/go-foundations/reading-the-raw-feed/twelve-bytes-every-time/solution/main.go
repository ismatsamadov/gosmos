package main

// sameBeacon reports whether all twelve bytes match.
func sameBeacon(a, b [12]byte) bool { return a == b }

// withSatellite returns a copy with byte 2, the satellite number, set to number.
func withSatellite(beacon [12]byte, number byte) [12]byte { beacon[2] = number; return beacon }
