package main

import "net/http"

// signedIn returns the session ID and the customer that the request's session cookie names, while that session lasts.
func (a *Account) signedIn(r *http.Request) (string, Customer, bool) {
	return "", Customer{}, false
}
