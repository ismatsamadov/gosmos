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
	Code    string `json:"error"`
	Message string `json:"message"`
}

// writeError answers with the status, code and message of the *APIError in err, wrapped or not, or with a 500 for any other error.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	if e, ok := errors.AsType[*APIError](err); ok {
		writeJSON(w, e.Status, ErrorBody{Code: e.Code, Message: e.Message})
		return
	}
	writeJSON(w, http.StatusInternalServerError, ErrorBody{Code: "internal", Message: "internal error"})
}

// invalidBooking returns the 422 for a booking that breaks a rule, with the rule's message.
func invalidBooking(msg string) *APIError {
	return &APIError{Status: http.StatusUnprocessableEntity, Code: "invalid_booking", Message: msg}
}
