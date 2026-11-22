package main

import "net/http"

// Account is the account page beside the API, where a customer's people sign in with a browser and see their own bookings.
type Account struct {
	Store     Store
	Customers map[string]Customer // by ID
	Sessions  *Sessions
}

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

// page shows the signed-in customer's own bookings and renews the browser's cookie; a browser without a session is sent to sign in.
func (a *Account) page(w http.ResponseWriter, r *http.Request) {
	id, c, ok := a.signedIn(r)
	if !ok {
		http.Redirect(w, r, "/account/sign-in", http.StatusSeeOther)
		return
	}
	list, err := a.Store.List(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	mine := []Booking{}
	for _, b := range list {
		if owns(c, b.SatID) {
			mine = append(mine, b)
		}
	}
	http.SetCookie(w, sessionCookie(id))
	renderPage(w, c, mine)
}
