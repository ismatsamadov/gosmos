package main

import "net/http"

// App is the assembled service. This lesson needs only the service and the handler that serves it.
type App struct {
	Service *Service
	Handler http.Handler // the mux inside every middleware
}
