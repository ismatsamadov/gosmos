package main

import "net/http"

// requestLine returns the first line req sends: its method, its path with any query, and its version.
func requestLine(req *http.Request) string {
	return req.Method + " " + req.URL.RequestURI() + " " + req.Proto
}
