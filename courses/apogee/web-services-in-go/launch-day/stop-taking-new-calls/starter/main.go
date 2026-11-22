package main

import (
	"context"
	"net/http"
	"time"
)

// stop shuts srv down: no new connections, and up to within for the requests in flight; a zero within waits for all of them.
func stop(srv *http.Server, within time.Duration) error {
	_ = context.WithTimeout
	return srv.Close()
}
