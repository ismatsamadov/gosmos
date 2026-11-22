package main

import (
	"io"
	"net/http"
	"os"
)

// sendFile answers with file, or with the part of it the request's Range asks for, so a broken download goes on where it stopped.
func sendFile(w http.ResponseWriter, r *http.Request, file *os.File) {
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, file) // TODO: answer the request's Range with http.ServeContent
}
