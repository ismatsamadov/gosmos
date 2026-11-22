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
	// TODO: mark the other six fields, 2 to 7, as Sat is marked
	fmt.Println(string(pic))
	fmt.Printf("%d bytes, aligned to %d\n", unsafe.Sizeof(g), unsafe.Alignof(g))
	fmt.Printf("%d candidates: %d bytes\n", candidates, candidates*unsafe.Sizeof(g))
}
