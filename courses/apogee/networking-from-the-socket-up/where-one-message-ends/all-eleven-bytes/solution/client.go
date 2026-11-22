package main

import "io"

// readExactly reads exactly n bytes from r, however the stream splits them.
func readExactly(r io.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}
