package main

import "io"

// captureText reads src while copying the consumed bytes to log.
func captureText(src io.Reader, log io.Writer) (string, error) {
	data, err := io.ReadAll(io.TeeReader(src, log))
	return string(data), err
}
