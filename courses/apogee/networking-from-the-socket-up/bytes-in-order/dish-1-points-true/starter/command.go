package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

// pointBody is a POINT's body: the azimuth, then the elevation, in hundredths of a degree.
func pointBody(azDeg, elDeg float64) ([]byte, error) {
	az := math.Round(azDeg * 100)
	el := math.Round(elDeg * 100)
	if az < 0 || az > 35999 {
		return nil, fmt.Errorf("azimuth %v is outside 0 to 359.99", azDeg)
	}
	if el < 0 || el > 9000 {
		return nil, fmt.Errorf("elevation %v is outside 0 to 90.00", elDeg)
	}
	body := binary.BigEndian.AppendUint16(nil, uint16(az))
	return binary.BigEndian.AppendUint16(body, uint16(int16(el))), nil
}

// command sends one command and reads its answer: nil for an ACK of its seq, the
// *NakError of a NAK, and an error naming the type and seq of anything else.
func (c *Client) command(f Frame) error {
	return c.send(f)
}
