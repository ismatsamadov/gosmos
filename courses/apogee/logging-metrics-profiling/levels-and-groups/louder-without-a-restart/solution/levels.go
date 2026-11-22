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
		if err := v.UnmarshalText(body); err != nil {
			http.Error(w, fmt.Sprintf("unknown level %q", body), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	return mux
}
