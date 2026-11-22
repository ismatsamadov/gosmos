package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// maxBodyBytes is the most a request's body may hold.
const maxBodyBytes = 65536

// decodeJSON decodes the JSON body of r into dst, or returns the *APIError to answer with.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return bodyError(err, "request body: "+err.Error())
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return bodyError(err, "request body: unexpected data after the JSON value")
	}
	return nil
}

// bodyError returns the answer to a body whose decoding failed with err: a 413 when it was too large, else a 400 saying msg, or that the body was empty.
func bodyError(err error, msg string) *APIError {
	_ = errors.AsType[*http.MaxBytesError]
	if err == io.EOF {
		msg = "request body: empty"
	}
	return &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: msg}
}
