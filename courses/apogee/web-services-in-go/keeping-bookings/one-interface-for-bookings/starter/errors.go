package main

import (
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

// invalidBooking returns the 422 for a booking that breaks a rule, with the rule's message.
func invalidBooking(msg string) *APIError {
	return &APIError{Status: http.StatusUnprocessableEntity, Code: "invalid_booking", Message: msg}
}
