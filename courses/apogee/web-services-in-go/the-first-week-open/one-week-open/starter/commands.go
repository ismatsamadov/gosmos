package main

import (
	"errors"
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

// uploadError returns the answer to a command form that couldn't be read: 415 for a body that isn't a form, 400 for a form without a file, 413 for one over the limit.
func uploadError(err error) error {
	switch {
	case errors.Is(err, http.ErrNotMultipart):
		return &APIError{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "send the command file as multipart/form-data"}
	case errors.Is(err, http.ErrMissingFile):
		return &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: `form has no "file" part`}
	}
	return bodyError(err, "request body: "+err.Error())
}
