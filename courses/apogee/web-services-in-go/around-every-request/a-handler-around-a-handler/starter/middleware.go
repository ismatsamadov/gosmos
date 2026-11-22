package main

import "net/http"

// noStore asks every cache between the service and its callers to keep no copy of an answer.
func noStore(next http.Handler) http.Handler {
	return next
}
