package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// Header is a request's fields: names in lower case, every value in order.
type Header map[string][]string

// Get returns the values of the field name, joined with ", ", or "" when it is absent.
func (h Header) Get(name string) string {
	return strings.Join(h[strings.ToLower(name)], ", ")
}

// Request is one request as the station server reads it.
type Request struct {
	Method  string // "PUT"
	Target  string // "/recordings/furrow-a/YST/18856"
	Header  Header
	Length  int64  // 0 for no body, -1 for chunked
	Body    []byte // read by the server before routing
	Trailer Header // fields after a chunked body
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

// readHeader reads header fields up to the blank line that ends them, refusing crooked lines
// and more than 32 fields.
func readHeader(br *bufio.Reader) (Header, error) {
	h := Header{}
	for fields := 0; ; fields++ {
		line, err := readLine(br)
		if errors.Is(err, errLineTooLong) {
			return nil, &RequestError{Status: 431, Msg: "header line over 1024 bytes"}
		}
		if err != nil {
			return nil, err
		}
		if line == "" {
			return h, nil
		}
		if fields == 32 {
			return nil, &RequestError{Status: 431, Msg: "more than 32 header fields"}
		}
		name, value, ok := strings.Cut(line, ":")
		if !ok || !isToken(name) || !isFieldValue(value) {
			return nil, &RequestError{Status: 400, Msg: "bad header line"}
		}
		name = strings.ToLower(name)
		h[name] = append(h[name], strings.Trim(value, " \t"))
	}
}

// isToken reports whether s is one or more of the characters a field name may hold.
func isToken(s string) bool {
	for _, r := range s {
		ok := r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' ||
			strings.ContainsRune("!#$%&'*+-.^_`|~", r)
		if !ok {
			return false
		}
	}
	return s != ""
}

// isFieldValue reports whether s is free of control characters other than the tab, as a
// field value must be. Bytes from 0x80 up are allowed: HTTP passes them through untouched.
func isFieldValue(s string) bool {
	for i := 0; i < len(s); i++ {
		if c := s[i]; c < ' ' && c != '\t' || c == 0x7f {
			return false
		}
	}
	return true
}

// contentLength decides a body's length from Content-Length, refusing every value two
// programs could read differently.
func contentLength(method string, h Header, max int64) (int64, error) {
	values := h["content-length"]
	switch {
	case len(values) > 1:
		return 0, &RequestError{Status: 400, Msg: "more than one content-length"}
	case len(values) == 0 && method == "PUT":
		return 0, &RequestError{Status: 411, Msg: "length required"}
	case len(values) == 0:
		return 0, nil
	}
	v := values[0]
	if v == "" || strings.Trim(v, "0123456789") != "" {
		return 0, &RequestError{Status: 400, Msg: fmt.Sprintf("bad content-length %q", v)}
	}
	if max <= 0 {
		max = math.MaxInt64
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil || n > max {
		return 0, &RequestError{Status: 413, Msg: fmt.Sprintf("body over %d bytes", max)}
	}
	if n > 0 && (method == "GET" || method == "HEAD") {
		return 0, &RequestError{Status: 400, Msg: method + " has no body"}
	}
	return n, nil
}

// readRequest reads a request's head from br, refusing a crooked one with its status, and
// decides its body's length with bodyLength. The body itself stays in br.
func readRequest(br *bufio.Reader, max int64) (*Request, error) {
	line, err := readLine(br)
	if errors.Is(err, errLineTooLong) {
		return nil, &RequestError{Status: 414, Msg: "request line over 1024 bytes"}
	}
	if err != nil {
		return nil, err
	}
	method, target, err := parseRequestLine(line)
	if err != nil {
		return nil, err
	}
	h, err := readHeader(br)
	if err != nil {
		return nil, err
	}
	switch len(h["host"]) {
	case 0:
		return nil, &RequestError{Status: 400, Msg: "missing host"}
	case 1:
	default:
		return nil, &RequestError{Status: 400, Msg: "more than one host"}
	}
	if method != "GET" && method != "HEAD" && method != "PUT" {
		return nil, &RequestError{Status: 501, Msg: fmt.Sprintf("method %q not implemented", method)}
	}
	n, err := bodyLength(method, h, max)
	if err != nil {
		return nil, err
	}
	return &Request{Method: method, Target: target, Header: h, Length: n}, nil
}
