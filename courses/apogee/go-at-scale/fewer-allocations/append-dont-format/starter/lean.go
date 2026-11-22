package main

import (
	"fmt"
	"time"
)

// appendContact appends c's dish loader line to dst and returns the
// result: AOS and LOS as 15:04:05, the dish and the pass ID, with a
// space between each two and a newline at the end.
func appendContact(dst []byte, c Contact) []byte {
	aos, los := c.AOS.Format(time.TimeOnly), c.LOS.Format(time.TimeOnly)
	line := fmt.Sprintf("%s %s %d %s\n", aos, los, c.Dish, c.ID)
	return append(dst, line...)
}
