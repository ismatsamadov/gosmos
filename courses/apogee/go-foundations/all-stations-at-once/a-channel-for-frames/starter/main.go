package main

import "fmt"

// poller sends the sequence numbers of three frames from Pipit-1's pass over Ystmark.
func poller(frames chan uint32) {
	for seq := uint32(4182); seq <= 4184; seq++ {
		frames <- seq
	}
}

func main() {
	frames := make(chan uint32)
	go poller(frames)
	seq := <-frames
	fmt.Println("archived", seq)
}
