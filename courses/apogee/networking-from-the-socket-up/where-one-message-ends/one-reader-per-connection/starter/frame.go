package main

import "encoding/binary"

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

// ByeError is the controller hanging up, and why.
type ByeError struct{ Reason string }

// Error says why the controller hung up: "controller closed: idle".
func (e *ByeError) Error() string {
	return "controller closed: " + e.Reason
}

// decodeFrame splits a frame's payload, the bytes after its length, into a Frame,
// without checking its CRC.
func decodeFrame(payload []byte) (Frame, error) {
	n := len(payload) - 4
	return Frame{Type: payload[0], Seq: binary.BigEndian.Uint32(payload[1:5]), Body: payload[5:n]}, nil
}
