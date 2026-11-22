package main

// APIError is an error that says exactly what the caller gets.
type APIError struct {
	Status  int    // 422
	Code    string // "invalid_booking"
	Message string // "los must be after aos"
}

// Error returns the message the caller reads.
func (e *APIError) Error() string { return e.Message }
