package main

import (
	"errors"
	"fmt"
	"net/http"
)

// Account is the account page beside the API, where a customer's people sign in with a browser and see their own bookings.
type Account struct {
	Store     Store
	Customers map[string]Customer // by ID
	Sessions  *Sessions
}

// routes registers the account page's handlers on a new mux.
func (a *Account) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /account/{$}", a.page)
	mux.HandleFunc("GET /account/sign-in", a.signInForm)
	mux.HandleFunc("POST /account/sign-in", a.signIn)
	mux.HandleFunc("POST /account/sign-out", a.signOut)
	mux.HandleFunc("POST /account/bookings/{id}/cancel", a.cancel)
	return mux
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

// signInForm shows the form that signs a browser in with its customer's API key.
func (a *Account) signInForm(w http.ResponseWriter, r *http.Request) {
	renderSignIn(w, http.StatusOK, "")
}

// signIn checks the API key the sign-in form sent, starts a session for its customer and hands the browser the session's cookie.
func (a *Account) signIn(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	c, ok := checkKey(a.Customers, r.PostFormValue("key"))
	if !ok {
		renderSignIn(w, http.StatusUnauthorized, "invalid api key")
		return
	}
	http.SetCookie(w, sessionCookie(a.Sessions.Start(c.ID)))
	http.Redirect(w, r, "/account/", http.StatusSeeOther)
}

// signOut ends the browser's session on the service and tells the browser to drop its cookie, then sends it to the sign-in form.
func (a *Account) signOut(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(cookieName); err == nil {
		a.Sessions.End(cookie.Value)
	}
	gone := sessionCookie("")
	gone.MaxAge = -1
	http.SetCookie(w, gone)
	http.Redirect(w, r, "/account/sign-in", http.StatusSeeOther)
}

// cancel cancels one of the signed-in customer's bookings; another customer's is not found, exactly like a missing one.
func (a *Account) cancel(w http.ResponseWriter, r *http.Request) {
	_, c, ok := a.signedIn(r)
	if !ok {
		http.Redirect(w, r, "/account/sign-in", http.StatusSeeOther)
		return
	}
	id := r.PathValue("id")
	b, err := a.Store.Get(r.Context(), id)
	if err == nil && !owns(c, b.SatID) {
		err = fmt.Errorf("booking %s: %w", id, ErrNotFound)
	}
	if err == nil {
		err = a.Store.Delete(r.Context(), id)
	}
	switch {
	case errors.Is(err, ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case err != nil:
		http.Error(w, "internal error", http.StatusInternalServerError)
	default:
		http.Redirect(w, r, "/account/", http.StatusSeeOther)
	}
}
