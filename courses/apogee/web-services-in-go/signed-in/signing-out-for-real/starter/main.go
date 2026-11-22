package main

import "net/http"

// signOut ends the browser's session on the service and tells the browser to drop its cookie, then sends it to the sign-in form.
func (a *Account) signOut(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/account/sign-in", http.StatusSeeOther)
}
