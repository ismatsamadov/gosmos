package main

import "crypto/rand"

// newSessionToken returns the token a browser keeps after logging in to the ops console.
func newSessionToken() string {
	return rand.Text()
}
