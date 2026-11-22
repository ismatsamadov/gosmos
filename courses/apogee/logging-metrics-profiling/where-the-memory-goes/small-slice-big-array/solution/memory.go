package main

import "bytes"

// tail returns a copy of the last n bytes of frames, or of all of them when there are fewer.
func tail(frames []byte, n int) []byte {
	return bytes.Clone(frames[max(len(frames)-n, 0):])
}
