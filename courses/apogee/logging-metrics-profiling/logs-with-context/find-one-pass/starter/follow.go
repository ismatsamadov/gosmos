package main

import (
	"bufio"
	"encoding/json"
	"io"
)

// follow reads JSON log lines from r and returns pass passID's lines, in the order read, and how many lines weren't JSON.
func follow(r io.Reader, passID string) (lines []Line, broken int, err error) {
	_, _ = bufio.NewScanner, json.Unmarshal // TODO: read r line by line and keep passID's lines
	return nil, 0, nil
}
