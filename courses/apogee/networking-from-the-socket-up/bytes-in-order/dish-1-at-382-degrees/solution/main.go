package main

import (
	"encoding/binary"
	"fmt"
)

// monday is Dish 1's status body at 09:14 on Monday 27 March, as its controller sent it.
const monday = "\x01\x02\x47\x95\x10\x4a\x02\x80\xff\x83\x00\x01\x00\x00\x00\x76"

func main() {
	body := []byte(monday)
	var order binary.ByteOrder = binary.BigEndian // as every field of the protocol is
	az := order.Uint16(body[2:4])                 // hundredths of a degree
	el := int16(order.Uint16(body[4:6]))          // hundredths of a degree
	wind := order.Uint16(body[6:8])               // tenths of a km/h
	fmt.Printf("az %.2f\n", float64(az)/100)
	fmt.Printf("el %.2f\n", float64(el)/100)
	fmt.Printf("wind %.1f\n", float64(wind)/10)
}
