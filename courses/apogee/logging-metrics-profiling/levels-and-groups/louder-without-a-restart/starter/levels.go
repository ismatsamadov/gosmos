package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// levelHandler serves the debug port's PUT /debug/level, which sets v from the body: debug, info, warn or error, in any case.
func levelHandler(v *slog.LevelVar) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /debug/level", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(io.LimitReader(r.Body, 64))
		if err != nil {
			http.Error(w, fmt.Sprintf("reading the level: %v", err), http.StatusBadRequest)
			return
		}
		_ = body // TODO: set v from the body, answering 400 for a level it doesn't know
		w.WriteHeader(http.StatusNoContent)
	})
	return mux
}
