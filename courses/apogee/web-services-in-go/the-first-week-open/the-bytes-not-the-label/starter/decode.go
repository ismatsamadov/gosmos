package main

import (
	"errors"
	"io"
	"net/http"
)

// maxBodyBytes is the most a request's body may hold.
const maxBodyBytes = 65536

// bodyError returns the answer to a body whose decoding failed with err: a 413 when it was too large, else a 400 saying msg, or that the body was empty.
func bodyError(err error, msg string) *APIError {
	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return &APIError{Status: http.StatusRequestEntityTooLarge, Code: "too_large", Message: "request body over 65536 bytes"}
	}
	if err == io.EOF {
		msg = "request body: empty"
	}
	return &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: msg}
}
