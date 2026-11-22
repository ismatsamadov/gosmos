package main

import (
	"fmt"
	"net/http"
)

// serveJSON writes the current view as JSON.
func (b *Board) serveJSON(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "TODO")
}
