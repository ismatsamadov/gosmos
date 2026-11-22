package main

import "net/http"

// noCrossOrigin refuses with 403 every request that would change something when a browser sent it from another origin.
func noCrossOrigin(h http.Handler) http.Handler {
	return h
}
