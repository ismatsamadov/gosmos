package main

import (
	"net/http"
	"strconv"
	"time"
)

// retryDelay returns how long a client should wait before trying again, as res's Retry-After says.
func retryDelay(res *http.Response) time.Duration {
	secs, err := strconv.Atoi(res.Header.Get("Retry-After"))
	switch {
	case err != nil || secs < 1:
		return time.Second
	case secs > 60:
		return time.Minute
	}
	return time.Duration(secs) * time.Second
}
