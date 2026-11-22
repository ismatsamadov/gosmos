package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

// APIError is an error that says exactly what the caller gets.
type APIError struct {
	Status  int    // 422
	Code    string // "invalid_booking"
	Message string // "los must be after aos"
}

// Error returns the message the caller reads.
func (e *APIError) Error() string { return e.Message }

// ErrorBody is the JSON body of every error the service writes.
type ErrorBody struct {
	Code      string `json:"error"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// writeError answers err as an ErrorBody: an *APIError with its own status, anything else as a 500.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	status, body := http.StatusInternalServerError, ErrorBody{Code: "internal", Message: "internal error"}
	if apiErr, ok := errors.AsType[*APIError](err); ok {
		status, body = apiErr.Status, ErrorBody{Code: apiErr.Code, Message: apiErr.Message}
	}
	body.RequestID = requestIDFrom(r.Context())
	writeJSON(w, status, body)
}

// writeJSON answers v as JSON with this status, or 500 in plain text when v can't be marshaled.
func writeJSON(w http.ResponseWriter, status int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(append(data, '\n'))
}
