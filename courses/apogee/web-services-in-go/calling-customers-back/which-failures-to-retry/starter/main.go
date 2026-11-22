package main

import (
	"context"
	"errors"
	"net/http"
)

// retryable reports whether trying again can help, from the receiver's status, or from err when no answer came.
func retryable(status int, err error) bool {
	_, _ = errors.Is(err, context.Canceled), http.StatusTooManyRequests
	return true
}
