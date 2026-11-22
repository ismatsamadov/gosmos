package main

import (
	"io"
	"net/http"
)

// readCommands returns the bytes of the form part named "file" in r's body, reading at most maxBodyBytes of the body.
func readCommands(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	// TODO: read only the part named "file", from a body capped at maxBodyBytes
	return io.ReadAll(r.Body)
}
