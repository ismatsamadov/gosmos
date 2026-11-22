package main

import "encoding/binary"

// decodeFrame splits a frame's payload, the bytes after its length, into a Frame,
// without checking its CRC.
func decodeFrame(payload []byte) (Frame, error) {
	n := len(payload) - 4
	return Frame{Type: payload[0], Seq: binary.BigEndian.Uint32(payload[1:5]), Body: payload[5:n]}, nil
}
