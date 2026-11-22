package main

import (
	"encoding/json"
	"net/http"
)

// writeJSON answers v as JSON with this status, or 500 in plain text when v can't be marshaled.
func writeJSON(w http.ResponseWriter, status int, v any) {
	json.NewEncoder(w).Encode(v)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}
