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
	if err := c.send(f); err != nil {
		return err
	}
	ans, err := c.recv()
	if err != nil {
		return err
	}
	switch {
	case ans.Type == TypeAck && ans.Seq == f.Seq:
		return nil
	case ans.Type == TypeNak:
		return nakError(ans)
	}
	return fmt.Errorf("answer to seq %d: type %#02x, seq %d", f.Seq, ans.Type, ans.Seq)
}

// jogBody is a JOG's body: the axis, 'A' or 'E', then the move in hundredths of a
// degree, from -5.00 to 5.00 and never 0.
func jogBody(axis byte, deltaDeg float64) ([]byte, error) {
	if axis != 'A' && axis != 'E' {
		return nil, fmt.Errorf("axis %q is not A or E", axis)
	}
	d := math.Round(deltaDeg * 100)
	if d == 0 || d < -500 || d > 500 {
		return nil, fmt.Errorf("jog %v is outside -5.00 to 5.00, or 0", deltaDeg)
	}
	return binary.BigEndian.AppendUint16([]byte{axis}, uint16(int16(d))), nil
}
