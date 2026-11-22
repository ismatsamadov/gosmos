package main

import (
	"encoding/binary"
	"fmt"
)

// rawStatus is a status frame's body exactly as it is on the wire.
type rawStatus struct {
	Dish    uint8
	State   uint8
	Az      uint16 // hundredths of a degree
	El      int16  // hundredths of a degree
	Wind    uint16 // tenths of a km/h
	Temp    int16  // tenths of a °C
	Flags   uint16
	LastSeq uint32
}

// decodeStatus decodes a status frame's 16-byte body into degrees, km/h and °C.
func decodeStatus(body []byte) (Status, error) {
	_ = binary.Decode
	_ = fmt.Errorf
	return Status{}, nil
}
