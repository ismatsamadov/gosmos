package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

// decodeBeacon reads the sequence number and battery millivolts from a Pipit beacon in hex.
func decodeBeacon(hexText string) (seq, battMV uint16, err error) {
	b, err := hex.DecodeString(hexText)
	if err != nil {
		return 0, 0, fmt.Errorf("beacon: %w", err)
	}
	if len(b) != 12 {
		return 0, 0, fmt.Errorf("beacon: %d bytes, want 12", len(b))
	}
	if string(b[:2]) != "AP" {
		return 0, 0, errors.New("beacon: sync bytes are not AP")
	}
	return binary.BigEndian.Uint16(b[3:5]), binary.BigEndian.Uint16(b[6:8]), nil
}

func main() {
	seq, battMV, err := decodeBeacon("4150011056001ee6ffd30082")
	fmt.Println(seq, battMV, err)
}
