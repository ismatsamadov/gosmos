package main

import (
	"strconv"
	"time"
)

// appendContact appends c's dish loader line to dst and returns the
// result: AOS and LOS as 15:04:05, the dish and the pass ID, with a
// space between each two and a newline at the end.
func appendContact(dst []byte, c Contact) []byte {
	dst = c.AOS.AppendFormat(dst, time.TimeOnly)
	dst = append(dst, ' ')
	dst = c.LOS.AppendFormat(dst, time.TimeOnly)
	dst = append(dst, ' ')
	dst = strconv.AppendInt(dst, int64(c.Dish), 10)
	dst = append(dst, ' ')
	dst = append(dst, c.ID...)
	return append(dst, '\n')
}
