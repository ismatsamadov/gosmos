package main

import "net/http"

// Account is the account page beside the API, where a customer's people sign in with a browser and see their own bookings.
type Account struct {
	Store     Store
	Customers map[string]Customer // by ID
	Sessions  *Sessions
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
