package main

import "io"

// readText reads all available bytes, returning partial text with any read error.
func readText(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	return string(data), err
}
