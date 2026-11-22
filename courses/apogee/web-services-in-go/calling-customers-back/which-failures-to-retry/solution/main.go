package main

import (
	"context"
	"errors"
	"net/http"
)

// retryable reports whether trying again can help, from the receiver's status, or from err when no answer came.
func retryable(status int, err error) bool {
	switch {
	case errors.Is(err, context.Canceled):
		return false
	case status == 0:
		return err != nil
	case status == http.StatusRequestTimeout, status == http.StatusTooManyRequests:
		return true
	default:
		return status >= 500 && status <= 599
	}
}
