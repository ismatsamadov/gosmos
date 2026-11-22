package main

import "time"

// The pass line's frame types: the byte after a frame's length.
const (
	TypeHello   byte = 0x01 // the rehearsal names its pass and dish: "thistle-2/PEL/02136 1"
	TypeTLM     byte = 0x20 // the station relays one telemetry line
	TypeUplink  byte = 0x30 // the rehearsal sends one command line; seq is the command's number
	TypeFile    byte = 0x40 // a file begins: its chunk count (uint16), then its name
	TypeChunk   byte = 0x41 // one chunk of that file: its number (uint16, from 0), then its data
	TypeEnd     byte = 0x42 // the file is all sent: the whole file's CRC-32 (uint32)
	TypeWelcome byte = 0x81 // the station takes the pass: the dish (uint8), then last (uint64)
	TypeAck     byte = 0x90 // the satellite ran the command whose number is seq
	TypeNak     byte = 0x91 // refused: a code (uint8), then the reason; seq is the refused frame's
	TypeBye     byte = 0xFF // the station hangs up, and why: "los" or "stowed"
)

// NAK codes: the first byte of a NAK's body.
const (
	NakUnknownPass byte = 1 // no such contact at this station and dish now
	NakStowed      byte = 2 // the station's dishes are stowed
	NakBadCRC      byte = 3 // a frame whose CRC doesn't match
	NakHelloFirst  byte = 4 // anything but HELLO before the call's first HELLO
	NakUnknownType byte = 5 // a type not in the table
	NakRefused     byte = 6 // the satellite refused the command: its own reason follows
)

// The pass line's limits.
const (
	MinFrameLen = 9                // a frame's length counts the bytes after it: 9 with no body
	MaxFrameLen = 1024             // and 1024 with the longest
	MaxBody     = MaxFrameLen - 9  // 1015 bytes
	MaxChunk    = MaxBody - 2      // a chunk's data: 1013 bytes
	IdleLimit   = 90 * time.Second // every read must bring at least one byte within this
)
