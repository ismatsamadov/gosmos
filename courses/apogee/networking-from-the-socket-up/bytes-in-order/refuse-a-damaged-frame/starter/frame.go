package main

import "errors"

// Frame is one controller-protocol frame, without its length and CRC.
type Frame struct {
	Type byte
	Seq  uint32 // a command's sequence number; 0 for requests and their answers
	Body []byte
}

var (
	ErrBadLength = errors.New("bad length") // wrapped: "bad length 65535"
	ErrBadCRC    = errors.New("bad crc")
)
