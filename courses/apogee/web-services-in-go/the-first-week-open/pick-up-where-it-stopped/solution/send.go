package main

import (
	"net/http"
	"os"
)

// sendFile answers with file, or with the part of it the request's Range asks for, so a broken download goes on where it stopped.
func sendFile(w http.ResponseWriter, r *http.Request, file *os.File) {
	info, err := file.Stat()
	if err != nil {
		writeError(w, r, err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}
