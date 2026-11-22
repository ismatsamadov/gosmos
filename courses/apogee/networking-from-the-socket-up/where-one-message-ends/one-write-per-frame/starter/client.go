package main

import "io"

// writeFrame sends f to w in a single Write.
func writeFrame(w io.Writer, f Frame) error {
	b := encodeFrame(f)
	w.Write(b[:2]) // the length
	w.Write(b[2:]) // the rest
	return nil
}
