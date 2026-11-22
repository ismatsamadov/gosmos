package main

import "fmt"

// poller sends the sequence numbers of Pipit-1's frames from one pass over Ystmark.
func poller(frames chan uint32) {
	for seq := uint32(4182); seq <= 4185; seq++ {
		frames <- seq
	}
	close(frames)
}

func main() {
	frames := make(chan uint32)
	go poller(frames)
	n := 0
	for seq := range frames {
		fmt.Println("archived", seq)
		n++
	}
	fmt.Printf("pass over: %d frames\n", n)
}
