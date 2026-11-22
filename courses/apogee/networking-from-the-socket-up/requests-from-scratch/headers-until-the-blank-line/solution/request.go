package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Header is a request's fields: names in lower case, every value in order.
type Header map[string][]string

// Get returns the values of the field name, joined with ", ", or "" when it is absent.
func (h Header) Get(name string) string {
	return strings.Join(h[strings.ToLower(name)], ", ")
}

// RequestError is a request the station server refuses, and its answer.
type RequestError struct {
	Status int    // 400
	Msg    string // "bad request line": the answer's body, without its newline
}

// Error returns the status and the message, as in "400 bad request line".
func (e *RequestError) Error() string {
	return strconv.Itoa(e.Status) + " " + e.Msg
}

// errLineTooLong is a line longer than the reader's 1024-byte buffer.
var errLineTooLong = errors.New("line too long")

// readLine reads one line, which must end in CR LF and hold no other CR, and returns it
// without its CR LF.
func readLine(br *bufio.Reader) (string, error) {
	line, err := br.ReadSlice('\n')
	switch {
	case errors.Is(err, bufio.ErrBufferFull):
		return "", errLineTooLong
	case err == io.EOF && len(line) > 0:
		return "", io.ErrUnexpectedEOF
	case err != nil:
		return "", err
	}
	s, ok := strings.CutSuffix(string(line), "\r\n")
	if !ok || strings.Contains(s, "\r") {
		return "", &RequestError{Status: 400, Msg: "bad line ending"}
	}
	return s, nil
}

// parseRequestLine splits a request line into its method and target, refusing any other
// shape and any version but HTTP/1.1.
func parseRequestLine(line string) (method, target string, err error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 || parts[0] == "" || !strings.HasPrefix(parts[1], "/") || !versionShaped(parts[2]) {
		return "", "", &RequestError{Status: 400, Msg: "bad request line"}
	}
	if parts[2] != "HTTP/1.1" {
		return "", "", &RequestError{Status: 505, Msg: fmt.Sprintf("version %q not supported", parts[2])}
	}
	return parts[0], parts[1], nil
}

// versionShaped reports whether v has the shape HTTP/<digit>.<digit>.
func versionShaped(v string) bool {
	digit := func(b byte) bool { return b >= '0' && b <= '9' }
	return len(v) == 8 && strings.HasPrefix(v, "HTTP/") && digit(v[5]) && v[6] == '.' && digit(v[7])
}

// readHeader reads header fields up to the blank line that ends them.
func readHeader(br *bufio.Reader) (Header, error) {
	h := Header{}
	for {
		line, err := readLine(br)
		if err != nil {
			return nil, err
		}
		if line == "" {
			return h, nil
		}
		name, value, _ := strings.Cut(line, ":")
		name = strings.ToLower(name)
		h[name] = append(h[name], strings.Trim(value, " \t"))
	}
}
