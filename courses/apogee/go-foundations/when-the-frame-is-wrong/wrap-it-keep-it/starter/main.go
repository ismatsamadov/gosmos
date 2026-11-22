package main

import (
	"errors"
	"fmt"
)

// ErrDuplicate means a frame's sequence number isn't newer than the last one stored.
var ErrDuplicate = errors.New("duplicate frame")

// accept returns ErrDuplicate when sequence number seq isn't newer than last.
func accept(seq, last uint32) error {
	if seq <= last {
		return ErrDuplicate
	}
	return nil
}

// ingest checks frame seq from satellite sat and says which frame an error belongs to.
func ingest(sat string, seq, last uint32) error {
	if err := accept(seq, last); err != nil {
		return fmt.Errorf("%s seq %d: %v", sat, seq, err)
	}
	return nil
}

func main() {
	err := ingest("thistle", 4182, 4182)
	fmt.Println(err)
	fmt.Println("duplicate:", errors.Is(err, ErrDuplicate))
}
