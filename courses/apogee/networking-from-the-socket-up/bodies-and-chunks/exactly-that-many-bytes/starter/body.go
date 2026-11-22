package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// readBody reads a body of exactly n bytes from br, refusing one that ends early.
func readBody(br *bufio.Reader, n int64) ([]byte, error) {
	var body []byte
	for int64(len(body)) < n {
		line, err := br.ReadBytes('\n')
		body = append(body, line...)
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("body ended after %d of %d bytes", len(body), n)
		}
		if err != nil {
			return nil, err
		}
	}
	return body, nil
}
