package main

import "fmt"

// decoderFault names a known internal failure from the older decoder.
type decoderFault struct{ reason string }

// runKnownDecoder reports a known fault, but re-panics any unexpected failure.
func runKnownDecoder(decode func() int) (value int, err error) {
	defer func() {
		if problem := recover(); problem != nil {
			err = fmt.Errorf("decoder fault: %v", problem)
		}
	}()
	return decode(), nil
}
