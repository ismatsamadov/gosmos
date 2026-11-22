package main

import (
	"io"
	"net/http"
)

// readCommands returns the bytes of the form part named "file" in r's body, reading at most maxBodyBytes of the body.
func readCommands(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	form, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}
	for {
		part, err := form.NextPart()
		if err == io.EOF {
			return nil, http.ErrMissingFile
		}
		if err != nil {
			return nil, err
		}
		if part.FormName() == "file" {
			return io.ReadAll(part)
		}
	}
}
