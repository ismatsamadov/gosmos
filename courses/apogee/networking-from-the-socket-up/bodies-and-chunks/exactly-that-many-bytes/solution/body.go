package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
