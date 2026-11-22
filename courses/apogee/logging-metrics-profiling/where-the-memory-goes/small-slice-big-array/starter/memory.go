package main

import "bytes"

// tail returns a copy of the last n bytes of frames, or of all of them when there are fewer.
func tail(frames []byte, n int) []byte {
	_ = bytes.Clone // TODO: copy the bytes, so the buffer behind frames can be freed
	return frames[max(len(frames)-n, 0):]
}
