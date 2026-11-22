package main

import "net/http"

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
