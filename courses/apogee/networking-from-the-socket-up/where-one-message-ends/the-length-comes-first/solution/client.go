package main

import (
	"bufio"
	"encoding/binary"
	"io"
)

// readPayload reads one frame and returns the bytes after its length: type, seq, body
// and CRC.
func readPayload(br *bufio.Reader) ([]byte, error) {
	var length [2]byte
	if _, err := io.ReadFull(br, length[:]); err != nil {
		return nil, err
	}
	payload := make([]byte, binary.BigEndian.Uint16(length[:]))
	if _, err := io.ReadFull(br, payload); err != nil {
		return nil, err
	}
	return payload, nil
}
