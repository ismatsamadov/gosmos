package main

import (
	"encoding/binary"
	"hash/crc32"
)

// encodeFrame returns f as it goes on the wire: its length, type, seq and body, then the
// CRC-32 of the type, seq and body.
func encodeFrame(f Frame) []byte {
	b := make([]byte, 0, 11+len(f.Body))
	b = binary.BigEndian.AppendUint16(b, uint16(9+len(f.Body)))
	b = append(b, f.Type)
	b = binary.BigEndian.AppendUint32(b, f.Seq)
	b = append(b, f.Body...)
	_ = crc32.ChecksumIEEE
	// TODO: append the CRC-32 of the type, seq and body instead of zero
	return binary.BigEndian.AppendUint32(b, 0)
}
