package main

import "io"

// writeFrame sends f to w in a single Write.
func writeFrame(w io.Writer, f Frame) error {
	_, err := w.Write(encodeFrame(f))
	return err
}
