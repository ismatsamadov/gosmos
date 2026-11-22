package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// readBody reads a body of exactly n bytes from br, refusing one that ends early.
func readBody(br *bufio.Reader, n int64) ([]byte, error) {
	body := make([]byte, n)
	got, err := io.ReadFull(br, body)
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, &RequestError{Status: 400, Msg: fmt.Sprintf("body ended after %d of %d bytes", got, n)}
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}

// chunkSize reads the size at the start of a chunk's line: hex digits, then anything after a
// semicolon, which it ignores.
func chunkSize(line string) (int64, error) {
	digits, _, _ := strings.Cut(line, ";")
	digits = strings.TrimRight(digits, " \t")
	n, err := strconv.ParseInt(digits, 16, 64)
	if err != nil || len(digits) > 16 || strings.Trim(digits, "0123456789abcdefABCDEF") != "" {
		return 0, &RequestError{Status: 400, Msg: fmt.Sprintf("bad chunk size %q", line)}
	}
	return n, nil
}
