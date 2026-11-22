package main

import (
	"errors"
	"os"
	"strconv"
)

// reasons are the words after each status code, the same words net/http uses.
var reasons = map[int]string{
	200: "OK",
	201: "Created",
	400: "Bad Request",
	404: "Not Found",
	405: "Method Not Allowed",
	408: "Request Timeout",
	409: "Conflict",
	411: "Length Required",
	413: "Request Entity Too Large",
	414: "Request URI Too Long",
	431: "Request Header Fields Too Large",
	501: "Not Implemented",
	505: "HTTP Version Not Supported",
}

// dateLayout is how the Date field writes a time, in UTC.
const dateLayout = "Mon, 02 Jan 2006 15:04:05 GMT"

// Field is one header field of an answer, written as given.
type Field struct{ Name, Value string }

// answer is what a route decides; serveConn writes it.
type answer struct {
	status int
	fields []Field // Content-Type, Content-Length and the rest, in order
	body   []byte  // nil for HEAD
	feed   bool    // GET /feed: serveConn streams the feed instead, then hangs up
}

// textAnswer is a plain-text answer: msg and a newline, with its type and length.
func textAnswer(status int, msg string) answer {
	body := []byte(msg + "\n")
	return answer{status: status, body: body, fields: []Field{
		{"Content-Type", "text/plain; charset=utf-8"},
		{"Content-Length", strconv.Itoa(len(body))},
	}}
}

// failure is the answer to a request that couldn't be read or kept: a RequestError's own
// status and message, or 408 when a deadline passed. Any other error has no answer, status
// 0, because the client has gone.
func failure(err error) answer {
	if e, ok := errors.AsType[*RequestError](err); ok {
		return textAnswer(e.Status, e.Msg)
	}
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return textAnswer(408, "request timeout")
	}
	return answer{}
}
