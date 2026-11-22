package main

import "io"

// countingWriter passes every write to w and keeps in n the bytes w took, from the file's size at the start.
type countingWriter struct {
	w io.Writer
	n int64
}

// Write writes p to w and adds to n the bytes w took, even when the write fails.
func (c *countingWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	return n, err
}
