package main

import "io"

// copyRecording transfers src to dst and returns its byte count and error.
func copyRecording(dst io.Writer, src io.Reader) (int64, error) { return io.Copy(dst, src) }
