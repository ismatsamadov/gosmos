package main

import (
	"mime"
	"net/http"
)

// isJSON reports whether the Content-Type in h says the body is JSON.
func isJSON(h http.Header) bool {
	mediaType, _, err := mime.ParseMediaType(h.Get("Content-Type"))
	return err == nil && mediaType == "application/json"
}
