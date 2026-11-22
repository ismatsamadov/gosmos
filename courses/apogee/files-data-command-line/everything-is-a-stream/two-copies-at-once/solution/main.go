package main

import "io"

// mirrorReport copies src to the screen and archive writers, in that order.
func mirrorReport(src io.Reader, screen, archive io.Writer) (int64, error) {
	return io.Copy(io.MultiWriter(screen, archive), src)
}
