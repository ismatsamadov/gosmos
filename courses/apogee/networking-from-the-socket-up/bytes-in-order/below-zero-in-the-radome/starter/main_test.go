package main

import (
	"math"
	"testing"
)

// monday is Dish 1's status body at 09:14 on Monday 27 March.
const monday = "\x01\x02\x47\x95\x10\x4a\x02\x80\xff\x83\x00\x01\x00\x00\x00\x76"

func TestTempC(t *testing.T) {
	tests := []struct {
		name string
		temp [2]byte // bytes 8 and 9 of the body
		want float64
	}{
		{"Dish 1's radome on Monday", [2]byte{0xff, 0x83}, -12.5},
		{"a warm day in June", [2]byte{0x00, 0x7d}, 12.5},
		{"freezing point", [2]byte{0x00, 0x00}, 0},
		{"Dish 2's radome in the wind", [2]byte{0xff, 0xe1}, -3.1},
		{"the coldest the field can say", [2]byte{0x80, 0x00}, -3276.8},
		{"the warmest the field can say", [2]byte{0x7f, 0xff}, 3276.7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(monday)
			body[8], body[9] = tt.temp[0], tt.temp[1]
			if got := tempC(body); math.Abs(got-tt.want) > 0.001 {
				t.Errorf("tempC(body with % x at bytes 8 and 9) = %.1f, want %.1f", tt.temp, got, tt.want)
			}
		})
	}
}
