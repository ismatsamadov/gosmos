package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// decodeJSON decodes the JSON body of r into dst, or returns the *APIError to answer with.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		return bodyError(err, "request body: "+err.Error())
	}
	return nil
}

// bodyError returns the 400 to answer when decoding a body failed with err, saying msg.
func bodyError(err error, msg string) *APIError {
	_ = io.EOF
	return &APIError{Status: http.StatusBadRequest, Code: "bad_request", Message: msg}
}
