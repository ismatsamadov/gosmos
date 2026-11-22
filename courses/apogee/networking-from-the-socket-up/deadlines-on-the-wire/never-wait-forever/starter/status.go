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
	if len(body) != 16 {
		return Status{}, fmt.Errorf("status body is %d bytes, want 16", len(body))
	}
	var raw rawStatus
	if _, err := binary.Decode(body, binary.BigEndian, &raw); err != nil {
		return Status{}, err
	}
	return Status{
		Dish:    raw.Dish,
		State:   DishState(raw.State),
		AzDeg:   float64(raw.Az) / 100,
		ElDeg:   float64(raw.El) / 100,
		WindKmh: float64(raw.Wind) / 10,
		TempC:   float64(raw.Temp) / 10,
		Flags:   raw.Flags,
		LastSeq: raw.LastSeq,
	}, nil
}
