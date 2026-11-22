package main

import (
	"bufio"
	"encoding/binary"
	"io"
)

// readPayload reads one frame and returns the bytes after its length: type, seq, body
// and CRC.
func readPayload(br *bufio.Reader) ([]byte, error) {
	_ = io.ReadFull
	_ = binary.BigEndian
	return nil, nil
}
