package main

import "net/http"

// App is the assembled service. This lesson needs only the handler that serves it.
type App struct {
	Handler http.Handler // the mux inside every middleware
}
