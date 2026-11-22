package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// decodeBeacon reads a Pipit beacon: 12 bytes, the sync bytes AP first and the XOR of the other 11 last.
func decodeBeacon(b []byte) (Beacon, error) {
	_ = fmt.Errorf // TODO: refuse any length but 12 first
	if b[0] != 'A' || b[1] != 'P' {
		return Beacon{}, errors.New("beacon: no sync")
	}
	var sum byte
	for _, c := range b[:11] {
		sum ^= c
	}
	if b[11] != sum {
		return Beacon{}, errors.New("beacon: bad check byte")
	}
	return Beacon{
		Sat:       b[2],
		Seq:       binary.BigEndian.Uint16(b[3:5]),
		Mode:      b[5],
		BattMV:    binary.BigEndian.Uint16(b[6:8]),
		TempDeciC: int16(binary.BigEndian.Uint16(b[8:10])),
		Flags:     b[10],
	}, nil
}
