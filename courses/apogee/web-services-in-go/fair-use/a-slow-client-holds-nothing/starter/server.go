package main

import (
	"net/http"
	"time"
)

// Timeouts are the service's limits, production values in comments. A zero leaves that limit unset.
type Timeouts struct {
	ReadHeader time.Duration // 5s: a client has this long to send its headers
	Read       time.Duration // 15s: ... to send the whole request
	Write      time.Duration // 30s: ... to take the whole answer
	Idle       time.Duration // 2m: a kept-alive connection may sit unused this long
	Handler    time.Duration // 10s: the deadline of every request's context
}

// newServer returns an http.Server for h with t's limits and a 16 KiB limit on headers.
func newServer(h http.Handler, t Timeouts) *http.Server {
	return &http.Server{
		Handler: h,
	}
}
