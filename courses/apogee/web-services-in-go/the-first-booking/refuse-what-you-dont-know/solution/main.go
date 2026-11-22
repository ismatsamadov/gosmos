package main

import (
	"encoding/json"
	"io"
	"net/http"
)

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

// bodyError returns the 400 to answer when decoding a body failed with err, saying msg.
func bodyError(err error, msg string) *APIError {
	if err == io.EOF {
		msg = "request body: empty"
	}
	return &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: msg}
}
