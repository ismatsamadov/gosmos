package main

import (
	"bufio"
	"fmt"
	"io"
	"time"
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

// statusLine is an answer's first line: the version, the code and its reason, then CR LF.
func statusLine(code int) string {
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n", code, reasons[code])
}

// writeResponse writes one answer to w: the status line, Date, the fields in order, a blank
// line and the body, through one buffered writer flushed once.
func writeResponse(w io.Writer, now time.Time, status int, fields []Field, body []byte) error {
	bw := bufio.NewWriter(w)
	bw.WriteString(statusLine(status))
	bw.WriteString("Date: " + now.UTC().Format(dateLayout) + "\r\n")
	for _, f := range fields {
		bw.WriteString(f.Name + ": " + f.Value + "\r\n")
	}
	bw.WriteString("\r\n")
	bw.Write(body)
	return bw.Flush()
}
