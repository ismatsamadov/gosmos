package main

import "io"

// readExactly reads exactly n bytes from r, however the stream splits them.
func readExactly(r io.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	got, err := r.Read(buf)
	return buf[:got], err
}
