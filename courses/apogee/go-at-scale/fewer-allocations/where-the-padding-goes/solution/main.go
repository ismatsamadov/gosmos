package main

import (
	"fmt"
	"unsafe"
)

// candidates is how many candidates take-over day's window holds.
const candidates = 613496

// mark writes c over size bytes of pic from offset: the bytes one field takes.
func mark(pic []byte, offset, size uintptr, c byte) {
	for i := offset; i < offset+size; i++ {
		pic[i] = c
	}
}

func main() {
	var g glance
	pic := make([]byte, unsafe.Sizeof(g))
	for i := range pic {
		pic[i] = '.'
	}
	mark(pic, unsafe.Offsetof(g.Sat), unsafe.Sizeof(g.Sat), '1')
	mark(pic, unsafe.Offsetof(g.Station), unsafe.Sizeof(g.Station), '2')
	mark(pic, unsafe.Offsetof(g.AOS), unsafe.Sizeof(g.AOS), '3')
	mark(pic, unsafe.Offsetof(g.LOS), unsafe.Sizeof(g.LOS), '4')
	mark(pic, unsafe.Offsetof(g.MaxElevDeg), unsafe.Sizeof(g.MaxElevDeg), '5')
	mark(pic, unsafe.Offsetof(g.Dish), unsafe.Sizeof(g.Dish), '6')
	mark(pic, unsafe.Offsetof(g.Booked), unsafe.Sizeof(g.Booked), '7')
	fmt.Println(string(pic))
	fmt.Printf("%d bytes, aligned to %d\n", unsafe.Sizeof(g), unsafe.Alignof(g))
	fmt.Printf("%d candidates: %d bytes\n", candidates, candidates*unsafe.Sizeof(g))
}
