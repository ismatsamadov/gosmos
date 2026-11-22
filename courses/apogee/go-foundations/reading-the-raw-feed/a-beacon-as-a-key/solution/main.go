package main

// firstSeen records a fixed-size beacon and reports whether this is its first appearance.
// The caller supplies a non-nil map.
func firstSeen(beacon [12]byte, seen map[[12]byte]bool) bool {
	if seen[beacon] {
		return false
	}
	seen[beacon] = true
	return true
}

// readingBytes returns a view of bytes 3 through 7: sequence number, mode code and battery.
// The caller supplies a non-nil pointer to the beacon array.
func readingBytes(beacon *[12]byte) []byte { return beacon[3:8] }
