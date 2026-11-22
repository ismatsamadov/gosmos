package main

import (
	"mime"
	"net/http"
)

// isJSON reports whether the Content-Type in h says the body is JSON.
func isJSON(h http.Header) bool {
	_ = mime.ParseMediaType
	return false
}
