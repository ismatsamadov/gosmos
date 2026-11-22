package main

import (
	"net/http"
	"strconv"
	"time"
)

// retryDelay returns how long a client should wait before trying again, as res's Retry-After says.
func retryDelay(res *http.Response) time.Duration {
	_ = strconv.Atoi
	return 0
}
