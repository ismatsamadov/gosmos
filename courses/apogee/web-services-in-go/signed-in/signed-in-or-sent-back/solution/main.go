package main

import "net/http"

// signedIn returns the session ID and the customer that the request's session cookie names, while that session lasts.
func (a *Account) signedIn(r *http.Request) (string, Customer, bool) {
	cookie, err := r.Cookie(cookieName)
	if err != nil { // http.ErrNoCookie: this browser has no session
		return "", Customer{}, false
	}
	id, ok := a.Sessions.Lookup(cookie.Value)
	c, known := a.Customers[id]
	if !ok || !known {
		return "", Customer{}, false
	}
	return cookie.Value, c, true
}
