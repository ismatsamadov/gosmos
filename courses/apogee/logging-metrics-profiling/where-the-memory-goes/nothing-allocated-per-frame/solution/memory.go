package main

import "strconv"

// appendPassID appends the pass ID "sat/station/orbit", the orbit padded to five digits, to buf.
func appendPassID(buf []byte, sat, station string, orbit int) []byte {
	buf = append(buf, sat...)
	buf = append(buf, '/')
	buf = append(buf, station...)
	buf = append(buf, '/')
	for p := 10000; p > 1 && orbit < p; p /= 10 {
		buf = append(buf, '0')
	}
	return strconv.AppendInt(buf, int64(orbit), 10)
}
