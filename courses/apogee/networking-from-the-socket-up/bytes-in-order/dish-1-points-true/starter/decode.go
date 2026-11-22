package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// decodeFrame turns a frame's payload, the bytes after its length, into a Frame, refusing
// it unless its CRC matches the rest.
func decodeFrame(payload []byte) (Frame, error) {
	if len(payload) < 9 {
		return Frame{}, fmt.Errorf("%w %d", ErrBadLength, len(payload))
	}
	n := len(payload) - 4
	if crc32.ChecksumIEEE(payload[:n]) != binary.BigEndian.Uint32(payload[n:]) {
		return Frame{}, ErrBadCRC
	}
	return Frame{Type: payload[0], Seq: binary.BigEndian.Uint32(payload[1:5]), Body: payload[5:n]}, nil
}
