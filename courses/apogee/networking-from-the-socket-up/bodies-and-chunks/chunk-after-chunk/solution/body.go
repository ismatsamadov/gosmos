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

// readChunked reads a chunked body from br: chunk after chunk until the zero chunk, then the
// trailer's fields. Its chunks may hold max bytes in all.
func readChunked(br *bufio.Reader, max int64) ([]byte, Header, error) {
	var body []byte
	for {
		line, err := readLine(br)
		if errors.Is(err, errLineTooLong) {
			return nil, nil, &RequestError{Status: 400, Msg: "chunk line over 1024 bytes"}
		}
		if err != nil {
			return nil, nil, err
		}
		n, err := chunkSize(line)
		if err != nil {
			return nil, nil, err
		}
		if n == 0 {
			trailer, err := readHeader(br)
			return body, trailer, err
		}
		chunk := make([]byte, n+2)
		if _, err := io.ReadFull(br, chunk); err != nil {
			return nil, nil, err
		}
		if string(chunk[n:]) != "\r\n" {
			return nil, nil, &RequestError{Status: 400, Msg: "bad chunk ending"}
		}
		body = append(body, chunk[:n]...)
	}
}
