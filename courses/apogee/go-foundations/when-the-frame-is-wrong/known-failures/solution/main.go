package main

import (
	"errors"
	"fmt"
)

// ErrDuplicate means a frame's sequence number isn't newer than the last one stored.
var ErrDuplicate = errors.New("duplicate frame")

// accept returns an error when sequence number seq isn't newer than last.
func accept(seq, last uint32) error {
	if seq <= last {
		return ErrDuplicate
	}
	return nil
}

// shouldPage reports whether err must wake the on-call engineer.
func shouldPage(err error) bool {
	return err != nil && !errors.Is(err, ErrDuplicate)
}

func main() {
	err := accept(4182, 4182)
	fmt.Println(err, "page:", shouldPage(err))
}
