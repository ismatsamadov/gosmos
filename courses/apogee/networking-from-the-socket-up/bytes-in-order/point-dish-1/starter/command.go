package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

// pointBody is a POINT's body: the azimuth, then the elevation, in hundredths of a degree.
func pointBody(azDeg, elDeg float64) ([]byte, error) {
	_ = math.Round
	_ = binary.BigEndian
	_ = fmt.Errorf
	return nil, nil
}
