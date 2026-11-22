package main

import (
	"encoding/binary"
	"errors"
	"fmt"
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

var (
	ErrBadLength = errors.New("bad length") // wrapped: "bad length 65535"
	ErrBadCRC    = errors.New("bad crc")
)

// NakError is a frame the controller refused.
type NakError struct {
	Seq    uint32
	Code   uint8
	Reason string // "wind lock"
}

// Error says which frame was refused and why: "seq 121 refused: wind lock".
func (e *NakError) Error() string {
	return fmt.Sprintf("seq %d refused: %s", e.Seq, e.Reason)
}

// ByeError is the controller hanging up, and why.
type ByeError struct{ Reason string }

// Error says why the controller hung up: "controller closed: idle".
func (e *ByeError) Error() string {
	return "controller closed: " + e.Reason
}

// DishState is what a dish is doing.
type DishState uint8

// The states a controller reports.
const (
	Idle DishState = iota
	Moving
	Tracking
	Stowed
	Fault
)

// String names the state as the controller does: "TRACKING".
func (s DishState) String() string {
	switch s {
	case Idle:
		return "IDLE"
	case Moving:
		return "MOVING"
	case Tracking:
		return "TRACKING"
	case Stowed:
		return "STOWED"
	case Fault:
		return "FAULT"
	}
	return fmt.Sprintf("STATE(%d)", uint8(s))
}

// The status frame's flags.
const (
	FlagHeater   uint16 = 1 << 0
	FlagWindLock uint16 = 1 << 1
	FlagFault    uint16 = 1 << 2
)

// Status is a decoded status frame: one dish at one moment.
type Status struct {
	Dish    uint8
	State   DishState
	AzDeg   float64 // 183.25
	ElDeg   float64 // 41.70
	WindKmh float64 // 64.0
	TempC   float64 // -12.5
	Flags   uint16
	LastSeq uint32
}

// Welcome is a decoded WELCOME frame.
type Welcome struct {
	Dish    uint8
	FwMajor uint8
	FwMinor uint8
	LastSeq uint32
}

// decodeWelcome decodes a WELCOME's 7-byte body.
func decodeWelcome(body []byte) (Welcome, error) {
	if len(body) != 7 {
		return Welcome{}, fmt.Errorf("welcome body is %d bytes, want 7", len(body))
	}
	return Welcome{Dish: body[0], FwMajor: body[1], FwMinor: body[2], LastSeq: binary.BigEndian.Uint32(body[3:])}, nil
}

// nakError turns a NAK into the error it carries.
func nakError(f Frame) *NakError {
	e := &NakError{Seq: f.Seq}
	if len(f.Body) > 0 {
		e.Code, e.Reason = f.Body[0], string(f.Body[1:])
	}
	return e
}
