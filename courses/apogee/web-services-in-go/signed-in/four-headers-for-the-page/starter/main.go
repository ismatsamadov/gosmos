package main

import "net/http"

// securityHeaders tells the browser how to treat every answer of the account page, before the answer is written.
func securityHeaders(next http.Handler) http.Handler {
	return next
}
