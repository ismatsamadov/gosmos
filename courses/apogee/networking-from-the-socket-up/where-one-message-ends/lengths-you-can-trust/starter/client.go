package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

// readPayload reads one frame and returns the bytes after its length: type, seq, body
// and CRC.
func readPayload(br *bufio.Reader) ([]byte, error) {
	var length [2]byte
	if _, err := io.ReadFull(br, length[:]); err != nil {
		return nil, err
	}
	n := int(binary.BigEndian.Uint16(length[:]))
	// TODO: refuse a length under 9 or over 1024 with ErrBadLength, wrapped with the length
	_ = fmt.Errorf
	payload := make([]byte, n)
	if _, err := io.ReadFull(br, payload); err != nil {
		return nil, err
	}
	return payload, nil
}
