package main

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

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
