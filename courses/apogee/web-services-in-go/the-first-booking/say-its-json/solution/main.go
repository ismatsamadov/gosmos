package main

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
)

// maxBodyBytes is the most a request's body may hold.
const maxBodyBytes = 65536

// decodeJSON decodes the JSON body of r into dst, strictly and within the limit, or returns the *APIError to answer with.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return &APIError{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "send the body as application/json"}
	}
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
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
	if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return &APIError{Status: http.StatusRequestEntityTooLarge, Code: "too_large", Message: "request body over 65536 bytes"}
	}
	if err == io.EOF {
		msg = "request body: empty"
	}
	return &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: msg}
}
