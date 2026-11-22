package main

import (
	"log"
	"net/http"
)

// firstWeek serves the week's new routes beside the API, all inside the service's middleware, with bookings counted against plans and webhooks kept outside.
func firstWeek(app *App, f *Files, q *Quota, l *log.Logger) http.Handler {
	s := app.Service
	// TODO: guard the customer's routes, count bookings with metered, and deliver webhooks with webhookClient
	mux := http.NewServeMux()
	mux.HandleFunc("POST /bookings/{id}/commands", f.upload)
	mux.HandleFunc("GET /bookings/{id}/file", f.download)
	mux.HandleFunc("POST /bookings/{id}/link", f.link)
	mux.HandleFunc("GET /files/{id}", f.shared)
	mux.Handle("GET /go", goNext(l))
	mux.Handle("/", s.routes())
	return chain(mux, requestID, logRequests(l, s.Now), recoverPanics(l), noStore, withDeadline(s.Timeouts.Handler))
}
