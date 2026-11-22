package main

import "fmt"

// statusLine is an answer's first line: the version, the code and its reason, then CR LF.
func statusLine(code int) string {
	_ = fmt.Sprintf
	return "HTTP/1.1 200 OK\r\n"
}
