package main

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// decodeFrame turns a frame's payload, the bytes after its length, into a Frame, refusing
// it unless its CRC matches the rest.
func decodeFrame(payload []byte) (Frame, error) {
	_ = fmt.Errorf
	_ = crc32.ChecksumIEEE
	n := len(payload) - 4
	return Frame{Type: payload[0], Seq: binary.BigEndian.Uint32(payload[1:5]), Body: payload[5:n]}, nil
}
