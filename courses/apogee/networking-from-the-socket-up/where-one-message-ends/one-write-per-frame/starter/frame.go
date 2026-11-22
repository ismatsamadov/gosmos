package main

import (
	"encoding/binary"
	"hash/crc32"
)

// Frame is one controller-protocol frame, without its length and CRC.
type Frame struct {
	Type byte
	Seq  uint32 // a command's sequence number; 0 for requests and their answers
	Body []byte
}

// Frame types.
const (
	TypeHello   byte = 0x01
	TypePing    byte = 0x02
	TypeQuery   byte = 0x03
	TypePoint   byte = 0x10
	TypeStow    byte = 0x11
	TypeUnstow  byte = 0x12
	TypeJog     byte = 0x13
	TypeWelcome byte = 0x81
	TypePong    byte = 0x82
	TypeStatus  byte = 0x83
	TypeAck     byte = 0x90
	TypeNak     byte = 0x91
	TypeBye     byte = 0xFF
)

// encodeFrame returns f as it goes on the wire: its length, type, seq and body, then the
// CRC-32 of the type, seq and body.
func encodeFrame(f Frame) []byte {
	b := make([]byte, 0, 11+len(f.Body))
	b = binary.BigEndian.AppendUint16(b, uint16(9+len(f.Body)))
	b = append(b, f.Type)
	b = binary.BigEndian.AppendUint32(b, f.Seq)
	b = append(b, f.Body...)
	return binary.BigEndian.AppendUint32(b, crc32.ChecksumIEEE(b[2:]))
}
