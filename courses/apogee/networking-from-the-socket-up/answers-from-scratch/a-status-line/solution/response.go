package main

import "fmt"

// statusLine is an answer's first line: the version, the code and its reason, then CR LF.
func statusLine(code int) string {
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n", code, reasons[code])
}
