package main

import "encoding/binary"

// tempC reads the temperature inside the radome, in °C, from a status frame's body:
// bytes 8 and 9, in tenths of a degree.
func tempC(body []byte) float64 {
	return float64(binary.BigEndian.Uint16(body[8:10])) / 10
}
