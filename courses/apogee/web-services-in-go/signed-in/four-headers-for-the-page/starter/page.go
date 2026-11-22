package main

import (
	"bytes"
	"html/template"
	"net/http"
)

// pageHTML is the account page: the customer's bookings, each with a Cancel button, then Sign out.
const pageHTML = `<!doctype html>
<html lang="en">
<meta charset="utf-8">
<title>Apogee bookings: {{.Name}}</title>
<h1>{{.Name}}</h1>
<ul>
{{range .Bookings}}<li>{{.ID}}: {{.SatID}} over {{.Station}}, {{.AOS.UTC.Format "2006-01-02 15:04:05"}} to {{.LOS.UTC.Format "15:04:05"}} UTC
<form method="post" action="/account/bookings/{{.ID}}/cancel"><button>Cancel {{.ID}}</button></form></li>
{{else}}<li>No bookings yet.</li>
{{end}}</ul>
<form method="post" action="/account/sign-out"><button>Sign out</button></form>
</html>
`

// signInHTML is the sign-in form, with a message above it when there is one.
const signInHTML = `<!doctype html>
<html lang="en">
<meta charset="utf-8">
<title>Apogee bookings: sign in</title>
<h1>Sign in</h1>
{{if .}}<p>{{.}}</p>
{{end}}<form method="post" action="/account/sign-in">
<label>API key <input name="key" type="password" autocomplete="off"></label>
<button>Sign in</button>
</form>
</html>
`

// render executes the template text with data into a buffer, then sends it as HTML with this status, or 500 if it fails.
func render(w http.ResponseWriter, status int, text string, data any) {
	var buf bytes.Buffer
	if err := template.Must(template.New("page").Parse(text)).Execute(&buf, data); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

// renderPage writes the account page for customer c and its bookings.
func renderPage(w http.ResponseWriter, c Customer, list []Booking) {
	render(w, http.StatusOK, pageHTML, struct {
		Name     string
		Bookings []Booking
	}{c.Name, list})
}

// renderSignIn writes the sign-in form with this status and message.
func renderSignIn(w http.ResponseWriter, status int, msg string) {
	render(w, status, signInHTML, msg)
}
