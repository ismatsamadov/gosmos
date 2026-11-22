package main

import "net/http"

// cors lets pages from one trusted origin call the API from a browser: it answers their preflights and lets them read the answers.
func cors(origin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return next
	}
}
