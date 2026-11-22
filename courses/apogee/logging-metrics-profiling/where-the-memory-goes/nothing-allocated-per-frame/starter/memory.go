package main

import (
	"fmt"
	"strconv"
)

// appendPassID appends the pass ID "sat/station/orbit", the orbit padded to five digits, to buf.
func appendPassID(buf []byte, sat, station string, orbit int) []byte {
	_ = strconv.AppendInt // TODO: append the ID's parts one by one, without fmt
	return append(buf, fmt.Sprintf("%s/%s/%05d", sat, station, orbit)...)
}
