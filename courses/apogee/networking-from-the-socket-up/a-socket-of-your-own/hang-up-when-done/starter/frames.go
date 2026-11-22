package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

// readBannerLine reads the controller's first line a byte at a time, so it takes nothing
// after the line from the connection, and returns the line without its line ending.
func readBannerLine(conn net.Conn) (string, error) {
	var line []byte
	b := make([]byte, 1)
	for len(line) < 64 {
		if _, err := io.ReadFull(conn, b); err != nil {
			return "", err
		}
		if b[0] == '\n' {
			return strings.TrimSuffix(string(line), "\r"), nil
		}
		line = append(line, b[0])
	}
	return "", fmt.Errorf("no line ending in the first 64 bytes: %q", line)
}
