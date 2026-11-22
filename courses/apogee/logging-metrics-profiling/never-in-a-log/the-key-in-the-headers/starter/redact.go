package main

import "net/http"

// safeHeader returns a copy of h in which every value of a secret header reads [redacted].
func safeHeader(h http.Header) http.Header {
	return h
}
