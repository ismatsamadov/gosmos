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

// writeChunk writes p as one chunk of a chunked answer and sends it at once. An empty p
// writes nothing, since a chunk of size 0 ends the answer.
func writeChunk(w *bufio.Writer, p []byte) error {
	if len(p) == 0 {
		return nil
	}
	fmt.Fprintf(w, "%x\r\n", len(p))
	w.Write(p)
	w.WriteString("\r\n")
	return w.Flush()
}

// endChunks writes the last chunk, of size 0, and the blank line after it, and sends them.
func endChunks(w *bufio.Writer) error {
	w.WriteString("0\r\n\r\n")
	return w.Flush()
}

// feed streams the feed's answer after its head: one chunk per line, a newline after each,
// until lines closes, then the last chunk.
func (s *Server) feed(bw *bufio.Writer, lines <-chan string) error {
	for line := range lines {
		if err := writeChunk(bw, []byte(line+"\n")); err != nil {
			return err
		}
	}
	return endChunks(bw)
}
