package main

import (
	"fmt"
	"time"
)

// passLength returns the length of a pass that lasts minutes plus seconds.
func passLength(minutes, seconds int) time.Duration {
	return time.Duration(minutes + seconds)
}

func main() {
	fmt.Printf("%v\n", passLength(9, 12))
	fmt.Printf("%v\n", passLength(0, 45))
	fmt.Printf("%v\n", passLength(12, 0))
}
