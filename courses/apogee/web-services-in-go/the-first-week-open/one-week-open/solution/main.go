package main

import (
	"log"
	"net/http"
)

// firstWeek serves the week's new routes beside the API, all inside the service's middleware, with bookings counted against plans and webhooks kept outside.
func firstWeek(app *App, f *Files, q *Quota, l *log.Logger) http.Handler {
	s := app.Service
	s.Store = metered{Store: s.Store, quota: q}
	app.Notifier.Sender.Client = webhookClient()
	mux := http.NewServeMux()
	mux.Handle("POST /bookings/{id}/commands", s.private(f.upload))
	mux.Handle("GET /bookings/{id}/file", s.private(f.download))
	mux.Handle("POST /bookings/{id}/link", s.private(f.link))
	mux.HandleFunc("GET /files/{id}", f.shared)
	mux.Handle("GET /go", goNext(l))
	mux.Handle("/", s.routes())
	return chain(mux, requestID, logRequests(l, s.Now), recoverPanics(l), noStore, withDeadline(s.Timeouts.Handler))
}
