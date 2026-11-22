package main

import (
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
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
			return nil, nil, cutOff(err)
		}
		n, err := chunkSize(line)
		if err != nil {
			return nil, nil, err
		}
		if n == 0 {
			trailer, err := readHeader(br)
			return body, trailer, cutOff(err)
		}
		if max > 0 && n > max-int64(len(body)) {
			return nil, nil, &RequestError{Status: 413, Msg: fmt.Sprintf("body over %d bytes", max)}
		}
		chunk := make([]byte, n+2)
		if _, err := io.ReadFull(br, chunk); err != nil {
			return nil, nil, cutOff(err)
		}
		if string(chunk[n:]) != "\r\n" {
			return nil, nil, &RequestError{Status: 400, Msg: "bad chunk ending"}
		}
		body = append(body, chunk[:n]...)
	}
}

// cutOff turns the stream's end into the refusal of a recording cut off before its end.
func cutOff(err error) error {
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return &RequestError{Status: 400, Msg: "recording cut off"}
	}
	return err
}

// bodyLength decides how a request's body is framed: -1 for chunked, otherwise its length
// from Content-Length. A request with both, or with a coding the server doesn't know, is
// refused.
func bodyLength(method string, h Header, max int64) (int64, error) {
	_, hasLength := h["content-length"]
	_, hasCoding := h["transfer-encoding"]
	if !hasCoding {
		return contentLength(method, h, max)
	}
	coding := h.Get("transfer-encoding")
	codings := strings.Split(coding, ",")
	switch {
	case hasLength:
		return 0, &RequestError{Status: 400, Msg: "both content-length and transfer-encoding"}
	case !strings.EqualFold(strings.Trim(codings[len(codings)-1], " \t"), "chunked"):
		return 0, &RequestError{Status: 400, Msg: fmt.Sprintf("bad transfer-encoding %q", coding)}
	case !strings.EqualFold(coding, "chunked"):
		return 0, &RequestError{Status: 501, Msg: fmt.Sprintf("transfer-encoding %q not implemented", coding)}
	case method == "GET" || method == "HEAD":
		return 0, &RequestError{Status: 400, Msg: method + " has no body"}
	}
	return -1, nil
}

// checkCRC checks body against Apogee-Crc32, taken from the trailer when it is there and
// from the header otherwise, and returns the body's CRC-32.
func checkCRC(body []byte, h, trailer Header) (uint32, error) {
	crc := crc32.ChecksumIEEE(body)
	sent := trailer.Get("apogee-crc32")
	if sent == "" {
		sent = h.Get("apogee-crc32")
	}
	if sent == "" {
		return crc, nil
	}
	want, err := strconv.ParseUint(sent, 16, 32)
	if len(sent) != 8 || err != nil {
		return 0, &RequestError{Status: 400, Msg: fmt.Sprintf("bad apogee-crc32 %q", sent)}
	}
	if uint32(want) != crc {
		return 0, &RequestError{Status: 400, Msg: fmt.Sprintf("crc mismatch: sent %08x, received %08x", want, crc)}
	}
	return crc, nil
}
