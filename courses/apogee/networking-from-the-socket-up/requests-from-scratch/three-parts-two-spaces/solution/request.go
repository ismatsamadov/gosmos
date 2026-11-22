package main

import (
	"fmt"
	"strconv"
	"strings"
)

// RequestError is a request the station server refuses, and its answer.
type RequestError struct {
	Status int    // 400
	Msg    string // "bad request line": the answer's body, without its newline
}

// Error returns the status and the message, as in "400 bad request line".
func (e *RequestError) Error() string {
	return strconv.Itoa(e.Status) + " " + e.Msg
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
