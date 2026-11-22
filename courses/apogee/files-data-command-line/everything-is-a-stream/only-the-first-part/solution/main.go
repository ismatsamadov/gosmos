package main

import "io"

// preview reads at most n bytes and leaves later bytes in the source.
func preview(src io.Reader, n int64) (string, error) {
	data, err := io.ReadAll(io.LimitReader(src, n))
	return string(data), err
}
