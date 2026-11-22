package main

import "encoding/binary"

// Beacon is a Pipit beacon, decoded from its 12 bytes.
type Beacon struct {
	Sat       uint8  // byte 2: 1 for Pipit-1, 2 for Pipit-2
	Seq       uint16 // bytes 3 and 4, high byte first
	Mode      uint8  // byte 5: the mode's code
	BattMV    uint16 // bytes 6 and 7: the battery in millivolts
	TempDeciC int16  // bytes 8 and 9: tenths of a degree C
	Flags     uint8  // byte 10
}

// encodeBeacon writes bc as the 12 bytes a Pipit sends: the sync bytes AP, its fields with the
// high byte first, then the XOR of bytes 0 to 10.
func encodeBeacon(bc Beacon) []byte {
	b := make([]byte, 12)
	b[0], b[1], b[2] = 'A', 'P', bc.Sat
	binary.BigEndian.PutUint16(b[3:5], bc.Seq)
	b[5] = bc.Mode
	binary.BigEndian.PutUint16(b[6:8], bc.BattMV)
	binary.BigEndian.PutUint16(b[8:10], uint16(bc.TempDeciC))
	b[10] = bc.Flags
	for _, c := range b[:11] {
		b[11] ^= c
	}
	return b
}
