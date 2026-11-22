package main

import "net/http"

// APIError is an error that says exactly what the caller gets.
type APIError struct {
	Status  int    // 422
	Code    string // "invalid_booking"
	Message string // "los must be after aos"
}

// Error returns the message the caller reads.
func (e *APIError) Error() string { return e.Message }

// invalidBooking returns the 422 for a booking that breaks a rule, with the rule's message.
func invalidBooking(msg string) *APIError {
	return &APIError{Status: http.StatusUnprocessableEntity, Code: "invalid_booking", Message: msg}
}
