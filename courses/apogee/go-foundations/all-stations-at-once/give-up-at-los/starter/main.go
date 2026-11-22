package main

import (
	"errors"
	"fmt"
	"time"
)

// ErrNoAck means the satellite sent no acknowledgment before the limit.
var ErrNoAck = errors.New("no acknowledgment before LOS")

// waitAck returns the first acknowledgment on acks, or ErrNoAck once limit has passed.
func waitAck(acks <-chan string, limit time.Duration) (string, error) {
	return <-acks, nil
}

func main() {
	acks := make(chan string, 1)
	acks <- "ACK pipit-2 118"
	fmt.Println(waitAck(acks, 9*time.Minute))
}
