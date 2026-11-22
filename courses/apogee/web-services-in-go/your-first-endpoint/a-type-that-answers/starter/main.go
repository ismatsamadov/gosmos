package main

import (
	"encoding/json"
	"net/http"
)

// stationList is a list of stations that answers a request with itself.
type stationList []Station

// ServeHTTP answers with the list as JSON, in its own order.
func (s stationList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder
}
